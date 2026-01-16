# sfs-cli

Command-line interface for [Semantic File Search (SFS)](https://github.com/vcnt/sfs-api).

Upload, search, and manage files using semantic search powered by vector embeddings.

## Features

- **Upload** - Index files for semantic search
- **Search** - Find content using natural language queries
- **List** - View all indexed files
- **Delete** - Remove files from the index
- **Download** - Retrieve stored files
- **Config** - Manage API connection settings
- **Daemon** - Background service for automatic file watching
- **Watch** - Auto-sync folders

## Installation

### From Source

```bash
go install github.com/vcnt/sfs-cli@latest
```

### Build from Repository

```bash
git clone https://github.com/vcnt/sfs-cli.git
cd sfs-cli
go build -o sfs
sudo mv sfs /usr/local/bin/
```

## Quick Start

### 1. Configure API Connection

```bash
# Set your SFS API URL
sfs config set api_url https://your-sfs-api.com

# Set your API key
sfs config set api_key your-secret-key

# Verify configuration
sfs config list
```

### 2. Upload a File

```bash
# Upload a file
sfs upload document.pdf

# Update an existing file
sfs upload --update document.pdf
```

### 3. Search

```bash
# Basic search
sfs search "machine learning algorithms"

# Advanced search with options
sfs search "deployment best practices" --limit 10 --threshold 0.7
```

### 4. List Files

```bash
# List all files
sfs list

# Filter by prefix
sfs list --prefix docs_
```

### 5. Download & Delete

```bash
# Download a file
sfs download document.pdf

# Download with custom output name
sfs download document.pdf ./my-doc.pdf

# Delete a file
sfs delete document.pdf
```

### 6. Daemon Management

The daemon runs in the background to enable automatic file watching.

```bash
# Enable daemon to start automatically on boot
sfs daemon enable

# Disable automatic startup
sfs daemon disable

# Manually start the daemon
sfs daemon start

# Stop the daemon
sfs daemon stop

# Restart the daemon
sfs daemon restart

# Check daemon status
sfs daemon status
```

## Configuration File

Configuration is stored in `~/.config/sfs/config.yaml`:

```yaml
api_url: https://your-api.com
api_key: your-secret-key
watch_dirs:
  - /home/user/documents
  - /home/user/projects
```

## Development

### Running Tests

```bash
go test ./...
```

### Building

```bash
go build -o sfs
```

## License

See [LICENSE](LICENSE) file.

## Related Projects

- [sfs-api](https://github.com/vcnt/sfs-api) - Backend API server
- [sfs-desktop](https://github.com/vcnt/sfs-desktop) - Desktop GUI client
