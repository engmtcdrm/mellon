package encrypt

import (
	"bytes"
	"crypto/rand"
	"crypto/sha512"
	"errors"
	"log/slog"
	"math/big"
	"net"
	"os"
	"os/user"

	"github.com/fernet/fernet-go"
)

// The code below is intentionally uncommented in an attempt to obfuscate and
// make it harder to understand. There are intentional redundant calls to try
// again, to obfuscate the code even further.

func hashSHA(data ...[]byte) []byte {
	var d2h []byte

	if len(data) == 0 {
		n, err := rand.Int(rand.Reader, big.NewInt(9901))
		if err != nil {
			slog.Error("An error occurred", "error", err)
			os.Exit(99)
		}
		n = n.Add(n, big.NewInt(100))

		d2h = make([]byte, n.Int64())
		_, err = rand.Read(d2h)
		if err != nil {
			slog.Error("An error occurred", "error", err)
			os.Exit(99)
		}
	} else {
		d2h = data[0]
	}

	h := sha512.New()

	h.Write(d2h)

	return h.Sum(nil)
}

var size = len(hashSHA())

func getRandEncrypt(s int) ([]byte, error) {
	d := make([]byte, s)
	_, err := rand.Read(d)
	if err != nil {
		return nil, err
	}

	var k fernet.Key
	if err := k.Generate(); err != nil {
		return nil, err
	}

	t, err := fernet.EncryptAndSign(d, &k)
	if err != nil {
		return nil, err
	}

	return t, nil
}

func saltValue(k fernet.Key, data []byte, hu []byte) ([]byte, error) {
	r, err := rand.Int(rand.Reader, big.NewInt(9001))
	if err != nil {
		return nil, err
	}
	r = r.Add(r, big.NewInt(1000))

	e1, err := getRandEncrypt(int(r.Int64()))
	if err != nil {
		return nil, err
	}

	e2, err := getRandEncrypt(size)
	if err != nil {
		return nil, err
	}

	e3, err := fernet.EncryptAndSign(hashSHA([]byte(hu)), &k)
	if err != nil {
		return nil, err
	}

	e4, err := getRandEncrypt(int(r.Int64()))
	if err != nil {
		return nil, err
	}

	return append(append(append(append(e1, e2...), data...), e3...), e4...), nil
}

func createReadKey(keyPath string, hu string) (fernet.Key, error) {
	var k fernet.Key

	if _, err := os.Stat(keyPath); os.IsNotExist(err) {
		slog.Debug("Key file does not exist. Creating a new key.")

		if err := k.Generate(); err != nil {
			return fernet.Key{}, err
		}

		sk, err := saltValue(k, []byte(k.Encode()), []byte(hu))
		if err != nil {
			return fernet.Key{}, err
		}

		if err = os.WriteFile(keyPath, sk, 0600); err != nil {
			return fernet.Key{}, err
		}

		return k, nil
	} else {
		slog.Debug("Key file exists. Reading key from file.")

		d, err := os.ReadFile(keyPath)
		if err != nil {
			return fernet.Key{}, err
		}

		e1, err := getRandEncrypt(size)
		if err != nil {
			return fernet.Key{}, err
		}

		e2, err := getRandEncrypt(size)
		if err != nil {
			return fernet.Key{}, err
		}

		if err := k.Generate(); err != nil {
			return fernet.Key{}, err
		}

		ld := (len(d) - (len(e1) + len(k.Encode()) + len(e2))) / 2
		d = d[ld : len(d)-ld]
		kd := d[len(e1) : len(e1)+len(k.Encode())]
		hue := d[len(e1)+len(k.Encode()):]

		k3, err := fernet.DecodeKeys(string(kd))
		if err != nil {
			return fernet.Key{}, err
		}

		hud := fernet.VerifyAndDecrypt(hue, 0, k3)

		if !bytes.Equal(hud, hashSHA([]byte(hu))) {
			return fernet.Key{}, errors.New("an error occurred during key verification")
		}

		return *k3[0], nil
	}
}

// Encryptor/Decryptor
//
// Into the depths we dive, where the secrets lie...
type Tomb struct {
	hu string
	k  fernet.Key
}

func NewTomb(keyPath string) (*Tomb, error) {
	h, err := net.LookupAddr("127.0.0.1")
	if err != nil {
		return nil, err
	}

	cu, err := user.Current()
	if err != nil {
		return nil, err
	}

	hu := h[0] + cu.Username

	k, err := createReadKey(keyPath, hu)
	if err != nil {
		return nil, err
	}

	return &Tomb{
		hu: hu,
		k:  k,
	}, nil
}

func (tomb *Tomb) CheckPerms(checkData []byte) bool {
	return bytes.Equal(checkData, hashSHA([]byte(tomb.hu)))
}

func (tomb *Tomb) Encrypt(msg []byte) ([]byte, error) {
	e1, err := getRandEncrypt(size)
	if err != nil {
		return nil, err
	}

	e2, err := fernet.EncryptAndSign(msg, &tomb.k)
	if err != nil {
		return nil, err
	}

	e3, err := fernet.EncryptAndSign(hashSHA([]byte(tomb.hu)), &tomb.k)
	if err != nil {
		return nil, err
	}

	enc4, err := getRandEncrypt(size)
	if err != nil {
		return nil, err
	}

	return append(append(append(e1, e2...), e3...), enc4...), nil
}

func (tomb *Tomb) Decrypt(data []byte) ([]byte, error) {
	e1, err := getRandEncrypt(size)
	if err != nil {
		return nil, err
	}

	e2, err := fernet.EncryptAndSign(hashSHA([]byte(tomb.hu)), &tomb.k)
	if err != nil {
		return nil, err
	}

	e3, err := getRandEncrypt(size)
	if err != nil {
		return nil, err
	}

	td := data[len(e1) : len(data)-(len(e2)+len(e3))]
	hud := data[len(e1)+len(td) : len(data)-len(e3)]

	k, err := fernet.DecodeKeys(tomb.k.Encode())
	if err != nil {
		return nil, err
	}

	if tomb.CheckPerms(fernet.VerifyAndDecrypt(hud, 0, k)) {
		msg := fernet.VerifyAndDecrypt(td, 0, k)

		if msg != nil {
			return msg, nil
		}
	}

	return nil, errors.New("an error occurred during decryption")
}
