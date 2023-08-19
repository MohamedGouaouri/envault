# envault
A CLI tool that injects vault secrets into your application process.

The goal of this CLI is to inject those secrets into the application's enviornment variables, removing the overhead of managing them in CI/CD enviornment for example.
The tool is very similar (same approach) to the `infisical CLI` that can
 be found [here](https://github.com/Infisical/infisical/tree/main/cli). The difference is that `infisical CLI` interacts with a dedicated `infisical` platform.

## How to use it?
Make sure that you have a vault instance running with secret engine in place. By design, the secret engine can be used as a environment for app secrets (dev, prod...).
To enable a secret engine in vault, you can do the following:
```bash
vault secrets enable -path=dev-secrets kv
```
And in order to interact with this secret engine (puting and getting data), you can do:
### Put
```bash
vault kv put dev-secrets/api_key sh.efefe65487sd4654b
```
### Get
```bash
vault kv get dev-secrets/api_key
```
For more details about the secret engines, you can checkout the official documentation [here](https://developer.hashicorp.com/vault/docs/secrets).

### envault
For now, envault has only one command which is the `run` command, which allows you to spin up application process with secrets injected in its environment.

Make sure that `VAULT_ADDR` and a valid `VAULT_TOKEN` variables are set, which point to `vault` api address and a valid `vault` token with valid access rights.

To spin a process, use the following
```bash
envault run --env=<env> --path=<secret-path> -- npm start
```
- `env` is the env to pull secrets from. It's the secret engine name in vault's terminology.
- `<secret-path>` is vault path to secrets data.

