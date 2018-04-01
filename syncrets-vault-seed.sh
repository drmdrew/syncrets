#!/bin/sh

for attempt in 1 2 3 4 5; do
    if nc -w 5 -z vault-a 8200; then
        echo "vault-a is responding"
        break
    else
        echo "vault-a is not responding. Sleeping..."
        sleep 1
    fi
done

vault auth-enable userpass
vault write auth/userpass/users/player1 password=pacman

vault write secret/foo value="bar"
vault write secret/foo/bar value="foobar"
vault write secret/gilbert value="sullivan"
vault write secret/it/was/the/best/of/times value="it was the worst of times"


