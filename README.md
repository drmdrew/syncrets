# syncrets
A utility for synchronizing secrets between systems (like vault) and formats (like ejson, eyaml)

This project is a *WORK IN PROGRESS*, so consider it useful for experimentation but not ready
for production use.

## vault Example

To recursively copy the secrets from the vault running on hostA to the vault running on hostB:

```
syncrets -r vault://hostA:8200/secrets/foo/ vault://hostB:8200/secrets/bar/
```

