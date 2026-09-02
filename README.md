# Generate Environment File Action ⚙️

A powerful and flexible GitHub Action that dynamically generates environment (`.env`) files and exports variables to the `$GITHUB_ENV` context directly from your repository's Secrets and Variables. 

This action is especially useful for CI/CD pipelines, Docker builds, dynamic Kubernetes (K8s) deployments, and other deployment workflows where you need to securely pass GitHub Secrets and Variables to your application without hardcoding them.

## Features

- **Dynamic `.env` Generation**: Automatically creates `.env`, `.env.local`, or any custom-named environment files based on your configuration.
- **System Environment Export**: Exports variables directly to the `$GITHUB_ENV` context so subsequent steps in your workflow can access them instantly.
- **Template Support**: Safely uses a `.env.example` file as a template to update existing keys.
- **Exclusion Lists**: Easily filter out specific sensitive keys (e.g., `DOCKER_PASSWORD`, `SSH_KEY`) that you don't want exposed in a file or the system environment.

## Who can use this?
This action is perfect for developers and DevOps engineers who:
- Build Docker images and need to pass `.env` files into their containers.
- Manage dynamic Kubernetes (K8s) deployments and need to securely inject variables into ConfigMaps or Secrets on the fly.
- Run automated tests that require various environment variables.
- Deploy applications to servers (e.g., Node.js, Python, Go, PHP) and need to generate production `.env` files dynamically on the fly.

---

## Usage Examples

### 1. Basic Usage (Export to GITHUB_ENV only)
By default, the action will read your secrets and variables, export them to `$GITHUB_ENV`, and will **not** generate any physical files.

```yaml
steps:
  - name: Load Secrets to Environment
    uses: alamin899/generate-env-action@v1
    with:
      secrets_context: ${{ toJson(secrets) }}
      vars_context: ${{ toJson(vars) }}
      
  - name: Test Environment Variables
    run: |
      echo "My secret is available: $MY_SECRET_KEY"
```

### 2. Generate `.env` files
If your application requires physical `.env` files, you can specify the filenames you want to generate.

```yaml
steps:
  - name: Generate .env and .env.local
    uses: alamin899/generate-env-action@v1
    with:
      secrets_context: ${{ toJson(secrets) }}
      vars_context: ${{ toJson(vars) }}
      generate_env_files: '.env, .env.local'

  - name: Verify file generation
    run: cat .env
```

### 3. Advanced Usage (Exclude Keys & File Only)
You might want to generate a `.env` file for your app but exclude sensitive pipeline credentials (like your Docker registry password). You can also disable the `GITHUB_ENV` export if you only want the physical file.

```yaml
steps:
  - name: Generate safe .env file
    uses: alamin899/generate-env-action@v1
    with:
      secrets_context: ${{ toJson(secrets) }}
      vars_context: ${{ toJson(vars) }}
      generate_env_files: '.env'
      exclude_keys: 'DOCKER_USERNAME,DOCKER_PASSWORD,GITHUB_TOKEN'
      is_set_process_env: 'false'
```

---

##  Inputs Reference

| Input | Required | Default | Description |
|-------|----------|---------|-------------|
| `secrets_context` | **Yes** | N/A | The JSON string of your GitHub secrets. Pass `${{ toJson(secrets) }}`. |
| `vars_context` | **Yes** | N/A | The JSON string of your GitHub variables. Pass `${{ toJson(vars) }}`. |
| `generate_env_files`| No | `''` (empty) | Comma-separated list of filenames to generate (e.g., `.env, .env.local`). If left empty, no files are created. |
| `is_set_process_env`| No | `'true'` | If `'true'`, exports all generated environment variables directly to `$GITHUB_ENV` for use in subsequent steps. |
| `exclude_keys` | No | `''` (empty) | Comma-separated list of additional keys to exclude from the generated files and the system environment. |

## 📝 License
This project is licensed under the MIT License.
