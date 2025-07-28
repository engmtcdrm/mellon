#!/bin/zsh

max="${1:-10}"

# First make sure secrets exist
./create_secrets.zsh "$max"

for i in {1..$max}; do
    secret_name="s${i}"
    go run ../. delete -s "$secret_name" -f

    echo -ne "\r\x1b[2KDeleted secret $i of $max"
done

echo ""
