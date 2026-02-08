# go-etl Contributing Guide

Thank you for your interest in the go-etl project! We welcome and appreciate contributions in various forms, including but not limited to bug reports, feature requests, code improvements, documentation improvements, and more.

## Table of Contents

- [Code of Conduct](#code-of-conduct)
- [Getting Started](#getting-started)
- [Development Environment Setup](#development-environment-setup)
- [Project Structure](#project-structure)
- [Developing New Plugins](#developing-new-plugins)
- [Code Standards](#code-standards)
- [Submitting Code](#submitting-code)
- [Testing](#testing)
- [Documentation](#documentation)
- [Getting Help](#getting-help)

## Code of Conduct

Please treat all contributors with respect and maintain friendly and professional communication. We aim to create an open and inclusive community, and any inappropriate behavior is unacceptable.

## Getting Started

### Reporting Bugs

If you find a bug, please report it through GitHub Issues. Include the following information:

- Problem description
- Steps to reproduce
- Expected behavior
- Actual behavior
- Environment information (operating system, Go version, etc.)
- Possible solutions

### Submitting Feature Requests

If you have new feature suggestions, feel free to discuss them via GitHub Issues. Please describe:

- The feature you want to implement
- Use case scenario
- Possible implementation ideas
- Any relevant reference materials

### Submitting Pull Requests

1. Fork this project
2. Create your feature branch: `git checkout -b feature/amazing-feature`
3. Commit your changes: `git commit -m 'Add some amazing feature'`
4. Push to your branch: `git push origin feature/amazing-feature`
5. Submit a Pull Request

## Development Environment Setup

### Requirements

- Go 1.20 or higher
- GCC 4.8 or higher (Linux)
- MinGW-w64 environment (Windows, GCC 7.2.0 or higher)

### Getting the Source Code

```bash
cd ${GO_PATH}/src
git clone https://github.com/Breeze0806/go-etl.git "github.com/Breeze0806/go-etl"
cd github.com/Breeze0806/go-etl
```

### Installing Dependencies

**Linux:**
```bash
make dependencies
```

**Windows:**
```bash
release.bat
```

### Building the Project

**Linux:**
```bash
make release
```

**Windows:**
```bash
release.bat
```

### Excluding DB2 Dependencies

If DB2 support is not needed, set the environment variable before building:

**Linux:**
```bash
export IGNORE_PACKAGES=db2
make dependencies
make release
```

**Windows:**
```bash
set IGNORE_PACKAGES=db2
release.bat
```

## Project Structure

go-etl project uses a Framework + plugin architecture and mainly includes the following modules:

```
go-etl/
├── datax/                    # Data synchronization framework
│   ├── plugin/
│   │   ├── reader/          # Reader plugins
│   │   └── writer/          # Writer plugins
│   └── ...
├── element/                  # Data types and type conversions
├── storage/                  # Storage modules
│   ├── database/            # Database integration
│   └── stream/              # Data stream processing
├── tools/                   # Tool collection
│   ├── datax/              # Build and release tools
│   └── license/            # License tool
└── ...
```

### Core Module Description

- **datax**: Offline data synchronization framework similar to Alibaba's DataX
- **element**: Data type definitions and type conversions
- **storage/database**: Database basic integration and dialect interfaces
- **storage/stream/file**: File parsing (CSV, Excel, etc.)
- **tools/build**: Plugin registration and code generation
- **tools/license**: Automatic license addition

## Developing New Plugins

go-etl supports extending data sources through plugins. Below is a guide for developing new plugins.

### Creating a Reader Plugin

1. Use the template generation tool to create the plugin framework:

```bash
cd tools/datax/plugin
go run main.go -t reader -p Mysql
```

2. Modify the generated files:
   - Update plugin information in `plugin.json`
   - Implement Job interface in `job.go`
   - Implement Task interface in `task.go`
   - Implement Reader interface in `reader.go`

3. Register the plugin:

```bash
go generate ./...
```

### Creating a Writer Plugin

1. Use the template generation tool to create the plugin framework:

```bash
cd tools/datax/plugin
go run main.go -t writer -p Mysql
```

2. Modify the generated files:
   - Update plugin information in `plugin.json`
   - Implement Job interface in `job.go`
   - Implement Task interface in `task.go`
   - Implement Writer interface in `writer.go`

3. Register the plugin:

```bash
go generate ./...
```

### Plugin Interface Description

#### Reader Plugin Interface

Reader plugins need to implement the following interfaces:

**Job Interface:**
```go
Init(ctx context.Context) (err error)
Destroy(ctx context.Context) (err error)
Split(ctx context.Context, number int) ([]*config.JSON, error)
Prepare(ctx context.Context) error
Post(ctx context.Context) error
```

**Task Interface:**
```go
Init(ctx context.Context) (err error)
Destroy(ctx context.Context) (err error)
StartRead(ctx context.Context, sender plugin.RecordSender) error
Prepare(ctx context.Context) error
Post(ctx context.Context) error
```

#### Writer Plugin Interface

Writer plugins need to implement the following interfaces:

**Job Interface:**
```go
Init(ctx context.Context) (err error)
Destroy(ctx context.Context) (err error)
Split(ctx context.Context, number int) ([]*config.JSON, error)
Prepare(ctx context.Context) error
Post(ctx context.Context) error
```

**Task Interface:**
```go
Init(ctx context.Context) (err error)
Destroy(ctx context.Context) (err error)
StartWrite(ctx context.Context, receiver plugin.RecordReceiver) error
Prepare(ctx context.Context) error
Post(ctx context.Context) error
SupportFailOver() bool
```

### Database Plugin Development

If implementing a relational database plugin, it is recommended to:

1. Refer to the [Database Storage Developer Guide](storage/database/README.md)
2. Implement the `Querier` interface (Reader) or `Execer` interface (Writer)
3. Use the `dbms.StartRead` or `dbms.StartWrite` function

### File Stream Plugin Development

If implementing a two-dimensional table file plugin (such as CSV, Excel), it is recommended to:

1. Refer to the [Two-dimensional Table File Stream Storage Developer Guide](storage/stream/file/README.md)
2. Implement file parsing and generation logic

### Plugin Configuration Design Principles

- Use camelCase naming
- Configuration items should be orthogonal with no overlapping functionality
- Use JSON types reasonably
- Follow conventions of similar plugins

## Code Standards

### Code Formatting

The project uses `gofmt` for code formatting:

```bash
gofmt -s -w yourfile.go
```

### Adding License

Before submitting code, please run the following command to automatically add a license:

```bash
go run tools/license/main.go
```

### Naming Conventions

- Package names: concise and meaningful
- Function names: camelCase
- Constant names: UPPER_CASE with underscores
- Variable names: camelCase

### Comment Conventions

- Public functions and types require comments
- Complex logic requires detailed comments
- Use English for comments

## Submitting Code

### Commit Message Guidelines

- Use English to describe the changes
- Concisely describe what was changed
- Include related Issue numbers (if any)

Example:
```
Add MySQL reader plugin support

Implement basic MySQL data reading functionality with batch fetch.

Fixes #123
```

### Pull Request Requirements

- Code must pass all tests
- Follow project code standards
- Include necessary documentation updates
- Clearly describe the purpose and content of the PR

## Testing

### Running Tests

```bash
go test ./...
```

### Test Coverage

The project encourages improving test coverage. New features should include corresponding unit tests.

### Performance Testing

For performance-related changes, please provide a performance test report including:

- Test environment (hardware, operating system, etc.)
- Test data characteristics
- Test parameter configuration
- Performance comparison data

## Documentation

### Code Comments

- Public APIs must have comments
- Complex logic requires detailed explanations
- Use English for comments

### README Documentation

New plugins need to include the following documentation:

1. **Quick Introduction**: Plugin functionality and use cases
2. **Implementation Principles**: Underlying implementation principles
3. **Configuration Instructions**: JSON configuration examples and parameter descriptions
4. **Type Conversion**: Data type conversion rules
5. **Performance Report**: Performance test data
6. **Constraints and Limitations**: Usage restrictions and precautions
7. **FAQ**: Frequently asked questions

### Updating Related Documentation

If your changes affect user usage, please update:

- `README_USER.md`: User manual
- `README_USER_zh-CN.md`: Chinese user manual
- Plugin documentation

## Getting Help

If you encounter problems during contribution, you can get help through the following methods:

- Check the [Project Documentation](README.md)
- Check the [Developer Documentation](datax/README.md)
- Check the [User Manual](README_USER.md)
- Submit GitHub Issue for discussion
- QQ Group: 185188648

## Acknowledgments

Thanks to all developers who have contributed to the go-etl project!
