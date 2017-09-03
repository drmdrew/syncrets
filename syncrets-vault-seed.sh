#!/bin/sh

vault auth-enable userpass
vault write auth/userpass/users/player1 password=pacman


