#!/bin/sh

vault auth-enable userpass
vault write auth/userpass/users/player1 password=pacman

vault write secret/foo value="bar"
vault write secret/foo/bar value="foobar"
vault write secret/gilbert value="sullivan"
vault write secret/it/was/the/best/of/times value="it was the worst of times"


