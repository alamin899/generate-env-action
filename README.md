# generate-env-action

This is a custom GitHub Action that dynamically generates a `.env` file from your GitHub repository's secrets and variables. It also supports using a `.env.example` file as a template.

## Usage

In any target project repository, reference your action inside your workflow (e.g., `.github/workflows/deploy.yml`):

```yaml
name: CI/CD Pipeline

on:
  push:
    branches: [ main ]

jobs:
  build:
    runs-on: ubuntu-latest
    steps:
      - name: Checkout project source code
        uses: actions/checkout@v4

      - name: Generate .env file
        uses: alamin899/generate-env-action@v1
        with:
          secrets_context: ${{ toJson(secrets) }}
          vars_context: ${{ toJson(vars) }}
          exclude_keys: 'SOME_EXTRA_SECRET,ANOTHER_KEY'

      - name: Verify generated .env file
        run: cat .env
```

## Inputs

- `secrets_context`: Required. The JSON string of the secrets context (e.g. `${{ toJson(secrets) }}`).
- `vars_context`: Required. The JSON string of the vars context (e.g. `${{ toJson(vars) }}`).
- `exclude_keys`: Optional. Comma-separated list of additional keys to exclude from the generated `.env` file.
- `is_set_process_env`: Optional. If `'true'`, exports all generated environment variables to `$GITHUB_ENV`. Default: `'true'`.
- `generate_env_files`: Optional. Comma-separated list of environment files to generate (e.g. `.env,.env.local`). If empty, no files are generated. Default: `''` (empty).
