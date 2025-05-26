# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Build/Test Commands
- Build: `just build`
- Run: `just run`
- Test all: `just test` 
- Test single file: `go test -v ./path/to/package -run TestName`
- Test with coverage: `just test-coverage`
- Format code: `just fmt`
- Lint code: `just lint`
- Install dependencies: `just deps`
- Update dependencies: `just update-deps`

## Code Style Guidelines
- **NO COMMENTS**: Do not write any comments in code - zero comments allowed
- **Architecture**: Follow clean architecture principles, organizing code by domain first, then technical concerns
- **Imports**: Standard library first, third-party next, internal packages last, with blank lines between groups
- **Formatting**: Use standard Go formatting (`go fmt`), tabs for indentation
- **Types**: PascalCase for exported types (`TaxClass`), camelCase for unexported types
- **Naming**: Package names use lowercase single words, exported functions use PascalCase, unexported use camelCase
- **Error Handling**: Return errors wrapped with context using `fmt.Errorf("failed to...: %w", err)`, check immediately
- **Testing**: Use table-driven tests with descriptive names, test both success and error cases, use mocks for external dependencies