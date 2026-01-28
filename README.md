# homebase-go-lib

A Go library for [your project description].

## Installation

```bash
go get github.com/parf/homebase-go-lib
```

## Usage

```go
import "github.com/parf/homebase-go-lib"

func main() {
    // Your code here
}
```

## Development

### Prerequisites

- Go 1.21 or higher
- Make (optional, for using Makefile commands)

### Building

```bash
make build
```

### Running Tests

```bash
make test
```

### Running Tests with Coverage

```bash
make test-coverage
```

### Formatting Code

```bash
make fmt
```

### Linting

```bash
make lint
```

## Project Structure

```
.
├── .github/          # GitHub Actions CI/CD workflows
├── docs/             # Additional documentation
├── examples/         # Example usage code
├── internal/         # Private packages (cannot be imported by external projects)
├── testdata/         # Test fixtures and data files
├── homebase.go       # Main library code
├── homebase_test.go  # Tests
├── go.mod            # Go module definition
├── Makefile          # Build automation
└── README.md         # This file
```

## Contributing

Contributions are welcome. Please open an issue or submit a pull request.

## License

MIT License - see LICENSE file for details.
