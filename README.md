# Redirect Service

A simple HTTP redirect service built with Go and Echo framework that allows you to configure URL redirects via command line arguments.

## Features

- Configurable redirects via command line flags
- Support for parameterized URLs (e.g., `:user`, `:id`)
- Query parameter preservation
- Environment variable support for port configuration
- Request logging

## Usage

### Basic Syntax

```bash
go run . -base <BASE_URL> [PATH1 REDIRECT1 PATH2 REDIRECT2 ...]
```

### Parameters

- `-base <URL>`: Base URL for redirects (required)
- `-port <port>`: Port to run the server on (optional, defaults to 8080 or PORT env var)
- Redirect mappings: Pairs of path patterns and redirect paths

### Examples

```bash
# Single redirect mapping
go run . -base "https://example.com" "/user/:user" "/u/:user"

# Multiple redirect mappings
go run . -base "https://example.com" \
  "/user/:user" "/u/:user" \
  "/user/:user/tags/:tag" "/u/:user/t/:tag"
# Custom port
go run . -base "https://example.com" -port "3000" "/path/:user" "/redirecthere/:user"
```

### Environment Variables

- `PORT`: Default port for the server (fallback if `-port` flag not provided)
