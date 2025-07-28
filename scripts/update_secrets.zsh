#!/bin/zsh

max="${1:-10}"

# First make sure secrets exist
./create_secrets.zsh "$max"

for i in {1..$max}; do
    secret_name="s${i}"
    new_password="newsecret${i}"
    tmpfile=$(mktemp)
    echo "$new_password" > "$tmpfile"
    go run ../. update -s "$secret_name" -f "$tmpfile" -c

    echo -ne "\r\x1b[2KUpdated secret $i of $max"
done

echo ""
