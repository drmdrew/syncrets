# syncrets
A utility for synchronizing secrets between systems like
[Hashicorp vault][VAULT] and formats like ejson/eyaml. Think of it like an
_rsync for secrets_. Secrets need to be handled carefully and syncrets can
help transfer, list, export, and otherwise manage secrets between systems
and formats.

Here is a simple example of using syncrets to copy secrets between two
vault servers running locally:
```
syncrets sync vault://localhost:8200/secrets/ vault://localhost:8201/secrets/
```
*NOTE*: This project is a *WORK IN PROGRESS*, so consider it useful for
experimentation but not ready for production use. Use at your own risk.

## syncrets config file
Outside of test scenarios, it isn't likely that two instances of vault would
be running on localhost.  To make working with multiple vaults easier,
syncrets supports a `~/.syncrets/syncrets.yml` configuration file, for
example:
```
vault:
    vault-a:
        url: "http://localhost:8200"
        auth:
            method: token
        token:
            file: ~/.syncrets/.vault-a-token
    vault-b:
        url: "http://localhost:8201"
        auth:
            method: token
        token:
            file: ~/.syncrets/.vault-b-token
```
Using a configuration file allows you to refer to servers using the name
(alias) present in their section of the configuration file.

For example, using the configuration above we can now rewrite the
previous syncrets example like so:
```
syncrets sync vault://vault-a/secrets/foo/ vault://vault-b/secrets/bar/
```
Using our configuration file syncrets will now know to reach `vault-a` using
`http://localhost:8200` and to reach `vault-b` using `http://localhost:8201`
which saves you from having to type out the full scheme, hostname, and port
when building URLs to pass to syncrets.

## syncrets commands
### auth
The `auth` command allows you to confirm that the authentication method being
used for a vault server is valid. If the authentication is invalid, the
syncrets `auth` command may prompt the user to reauthenticate using the
authentication method configured for the server.

### list
To recursively list the secrets (just the keys, no values) of a vault server
running on localhost you can use the `list` command:
```
syncrets list vault://localhost:8200/secrets/
```

### sync
To recursively copy the secrets between two vault servers running on localhost
you can use the `sync` command:
```
syncrets sync vault://localhost:8200/secrets/foo/ vault://localhost:8201/secrets/bar/
```

### rm
To recursively remove secrets of a vault server running on localhost you can
use the `rm` command:
```
syncrets rm vault://localhost:8200/secrets/
```
*CAUTION*: Use the `rm` command _carefully_, it can be a potent footgun.


[VAULT]: https://www.vaultproject.io/
