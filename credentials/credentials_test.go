package credentials

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIsValidName(t *testing.T) {
	assert.True(t, IsValidName("valid_name"))
	assert.True(t, IsValidName("valid-name"))
	assert.True(t, IsValidName("validname123"))
	assert.False(t, IsValidName("invalid name"))
	assert.False(t, IsValidName("invalid/name"))
	assert.False(t, IsValidName("invalid.name"))
}

func TestIsExists(t *testing.T) {
	// Create test data
	os.MkdirAll("./testdata", os.ModePerm)
	defer func() {
		if removeErr := os.RemoveAll("./testdata"); removeErr != nil {
			t.Errorf("Failed to remove testdata directory: %v", removeErr)
		}
	}()

	credFile := "./testdata/test.cred"
	os.WriteFile(credFile, []byte("test"), os.ModePerm)

	assert.True(t, IsExists(credFile))
	assert.False(t, IsExists("./testdata/nonexistent.cred"))
}
