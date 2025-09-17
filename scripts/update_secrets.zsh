#!/bin/zsh

max="${1:-10}"

# First make sure secrets exist
./create_secrets.zsh "$max"

for i in {1..$max}; do
    secret_name="s${i}"
    new_secret="newsecret${i}${i}"
    tmpfile=$(mktemp)
    echo "$new_secret" > "$tmpfile"
    old_secret=$(go run ../. view -s "$secret_name")
    go run ../. update -s "$secret_name" -f "$tmpfile" -c
    new_secret=$(go run ../. view -s "$secret_name")

    echo -e "\r\x1b[2KUpdated secret $i of $max; old: $old_secret, new: $new_secret"
done

echo ""
