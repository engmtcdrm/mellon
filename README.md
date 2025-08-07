# mellon

[![Go](https://github.com/engmtcdrm/mellon/actions/workflows/build.yml/badge.svg)](https://github.com/engmtcdrm/mellon/actions/workflows/build.yml)
[![Release](https://img.shields.io/github/v/release/engmtcdrm/mellon.svg?label=Latest%20Release)](https://github.com/engmtcdrm/mellon/releases/latest)

> *"Speak friend and enter"* - A lightweight CLI tool for securing and obtaining secrets

## What is mellon?

**mellon** is a secure, lightweight command-line tool for managing secrets locally. Named after the Elvish word for "friend" from The Lord of the Rings, mellon provides an easy way to encrypt, store, and retrieve sensitive information like passwords, API keys, tokens, and other confidential data on your local machine.

## Key Features

- 🔐 **Strong Encryption**: Uses industry-standard encryption to protect your secrets
- 💾 **Local Storage**: All secrets are stored locally on your machine - no cloud dependencies
- 🎯 **Simple Interface**: Interactive prompts or command-line flags for easy usage
- 📁 **File-based Input**: Import secrets directly from files with optional cleanup
- 🔍 **Quick Access**: List, view, update, and delete secrets with simple commands
- 🛡️ **Secure Permissions**: Automatically sets secure file permissions (0600 for secrets, 0700 for directories)

## Installation

Download the latest release from the [releases page](https://github.com/engmtcdrm/mellon/releases/latest) or build from source:

```bash
go install github.com/engmtcdrm/mellon@latest
```

## Quick Start

### Create your first secret
```bash
# Interactive mode - you'll be prompted for the secret and name
mellon create

# Or specify directly via flags
mellon create -s "my-api-key" -f ./secret.txt

# Create and cleanup the source file
mellon create -s "database-password" -f ./db-pass.txt --cleanup
```

### View a secret
```bash
# Interactive selection
mellon view

# View specific secret
mellon view -s "my-api-key"

# Save decrypted secret to file
mellon view -s "my-api-key" -o ./decrypted-key.txt
```

### List all secrets
```bash
# Detailed list with metadata
mellon list

# Simple name-only list
mellon list --print
```

### Update a secret
```bash
# Interactive update
mellon update

# Update from file
mellon update -s "my-api-key" -f ./new-secret.txt
```

### Delete secrets
```bash
# Interactive deletion
mellon delete

# Delete specific secret
mellon delete -s "my-api-key"

# Force delete without confirmation
mellon delete -s "my-api-key" --force

# Delete all secrets (use with caution!)
mellon delete --all
```

## Usage Examples

### Managing API Keys
```bash
# Store a GitHub token
echo "ghp_xxxxxxxxxxxxxxxxxxxx" > github-token.txt
mellon create -s "github-token" -f github-token.txt --cleanup

# Use the token in a script
API_TOKEN=$(mellon view -s "github-token")
curl -H "Authorization: token $API_TOKEN" https://api.github.com/user
```

### Database Credentials
```bash
# Store database password
mellon create -s "prod-db-password"
# Enter password interactively when prompted

# Retrieve for connection
DB_PASS=$(mellon view -s "prod-db-password")
mysql -u admin -p"$DB_PASS" production_db
```

### SSH Keys and Certificates
```bash
# Store SSH private key
mellon create -s "deploy-key" -f ~/.ssh/deploy_key

# Restore when needed
mellon view -s "deploy-key" -o ~/.ssh/restored_key
chmod 600 ~/.ssh/restored_key
```

## Command Reference

```bash
Usage:
  mellon [command]

Available Commands:
  create      Create a secret
  delete      Delete a secret
  help        Help about any command
  list        List available secrets
  update      Update a secret
  view        View a secret

Flags:
  -h, --help      help for mellon
  -v, --version   version for mellon

Use "mellon [command] --help" for more information about a command.
```

### Command Details

| Command | Description | Key Flags |
|---------|-------------|-----------|
| `create` | Encrypt and store a new secret | `-s` (secret name), `-f` (input file), `-c` (cleanup file) |
| `view` | Decrypt and display a secret | `-s` (secret name), `-o` (output file) |
| `update` | Modify an existing secret | `-s` (secret name), `-f` (input file), `-c` (cleanup file) |
| `list` | Show all stored secrets | `--print` (names only) |
| `delete` | Remove secrets | `-s` (secret name), `--force` (skip confirmation), `--all` (delete all) |

## Security

- **Encryption**: mellon uses strong encryption algorithms to protect your secrets
- **File Permissions**: Automatically sets restrictive permissions on secret files (0600) and directories (0700)
- **Local Only**: All data stays on your local machine - no network requests or cloud storage
- **Memory Safety**: Sensitive data is cleared from memory after use where possible

## Storage Location

By default, mellon stores encrypted secrets in:
- **Linux/macOS**: `~/.mellon/.thurin/`
- **Windows**: `%USERPROFILE%\.mellon\.thurin\`

The encryption key is stored separately from the secrets for additional security.

## Contributing

Contributions are welcome! Please feel free to submit issues and pull requests to the [GitHub repository](https://github.com/engmtcdrm/mellon).

## License

This project is licensed under the terms specified in the [LICENSE](LICENSE) file.
