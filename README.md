# Article Validator

CLI tool to validate bibliographic references via academic APIs.

## Features
- Parsing references from raw text
- Validation via Semantic Scholar API
- Concurrent API checks

## Installation
```bash
go install github.com/Lontor/article-validator/cmd/cli@latest
```

## Usage
```bash
validator "Einstein A. Relativity: The Special and General Theory. 1916"
```

## Development
```bash
# Run tests with coverage
go test -cover ./...

# Build binary
go build -o validator ./cmd/cli
```

## License
MIT