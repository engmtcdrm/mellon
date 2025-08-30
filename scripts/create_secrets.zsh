#!/bin/zsh

max="${1:-10}"

go run ../. delete --all -f

for i in {1..$max}; do
    secret_name="s${i}"
    secret="supersecret${i}"
    tmpfile=$(mktemp)
    echo "$secret" > "$tmpfile"
    go run ../. create -s "$secret_name" -f "$tmpfile" -c

    echo -ne "\r\x1b[2KCreated secret $i of $max"
done

echo ""
