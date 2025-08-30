#!/bin/zsh

max="${1:-10}"

# First make sure secrets exist
./create_secrets.zsh "$max"

for i in {1..$max}; do
    secret_name="s${i}"
    echo "$secret_name secret is " $(go run ../. view -s "$secret_name")
done
