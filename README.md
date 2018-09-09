# syncrets
*WIP*: This project is a *WORK IN PROGRESS*, so consider it useful for
experimentation but not ready for production use. Use at your own risk
but if you *do* use it I would love to hear what you think so please
log issues for anything you would like to see fixed/improved.

syncrets is a little utility for synchronizing secrets between systems like
[Hashicorp vault][VAULT] and formats like [ejson][EJSON]. Think of it like an
_rsync for secrets_. Secrets need to be handled carefully and syncrets can
help transfer, list, export, and otherwise manage secrets between systems
and formats. The name `syncrets` is a portmanteau of `secrets` and `sync` ...
obligatory [xkcd][XKCD-739].

Here is a simple example of using syncrets to copy secrets between two
vault servers running locally:
```
syncrets sync vault://vault-a/secrets/ vault://vault-b/secrets/
```

## syncrets config file

To faciliate working with multiple vaults, syncrets looks for a `syncrets.yml`
configuration file in the working directory as well as `~/.syncrets/syncrets.yml`.
Here is an example:

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
(alias) present in their section of the configuration file, so you can
refer to `vault://vault-a/secrets` rather than `http://localhost:8200/secrets`.

This example configuration file configures syncrets to reach `vault-a` using
`http://localhost:8200` and to reach `vault-b` using `http://localhost:8201`
which saves you from having to type out the full scheme, hostname, and port
when building URLs to pass to syncrets. The configuration also tells syncrets to
load vault auth tokens from file (assuming that these tokens have been obtained
previously).

## syncrets ejson

syncrets can directly `sync` secrets between two vault servers but can also
be used to `sync` secrets to a local file (preferrably in ejson format ...
these are _secrets_ after all).

If the source or target of a syncrets `sync` ends with `.ejson` then
syncrets will use the `ejson` configuration section of `syncrets.yml` to
configure the default encryption public key to use:
```
ejson:
    public_key:   a9d52487a1232e5c292a9680f4a44a84ea302ba05ff12d2e9d11662d20fc0139
```

For both encryption and decryption syncrets assumes that the ejson `EJSON_KEYDIR`
environment has been set if the ejson keys are not present in their default location.

*Example*:
```
syncrets sync vault://vault-a/secret/ ./secrets.ejson
```

Note: syncrets will write _unencrypted_ secrets to files ending with `.json` but
this regular JSON format is included primarily for testing/debugging purposes and
shouldn't be used for anything that is sensitive if the underlying filesystem isn't
trustworthy.

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
[EJSON]: https://github.com/Shopify/ejson
[XKCD-739]: https://xkcd.com/739/
