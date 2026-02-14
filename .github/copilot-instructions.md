# Copilot Instructions for endobit.io/table

## Commit Conventions

**Always use Conventional Commits format:**
- `fix:` for bug fixes
- `feat:` for new features
- `docs:` for documentation changes
- `chore:` for maintenance tasks
- `refactor:` for code refactoring
- `test:` for test additions/changes

Example: `fix: correct WithColor nil check and package comment typo`

## Development Workflow

**Before committing any code changes:**
1. Run `make test` to ensure all tests pass
2. Run `make lint` to check for linting issues
3. Fix any linter errors before committing

The project uses golangci-lint with 50+ enabled linters including:
- Style checkers (wsl_v5, whitespace, nlreturn)
- Code quality (gocritic, revive, staticcheck)
- Security (gosec)
- Performance (mirror, perfsprint)
- Best practices (errchkjson, errorlint, nilerr)

Common linter fixes:
- Use `strings.Contains(str, "text")` instead of `bytes.Contains([]byte(str), []byte("text"))`
- Always handle errors: `_ = func()` or `if err := func(); err != nil`
- Add blank lines after variable declarations and between logical blocks
- Replace simple lambdas with direct function references when possible

## Build, Test, and Lint Commands

This project uses a custom build system via `endobit.io/builder` with Makefile targets:

```bash
# Run tests (with coverage)
make test                      # Full test suite
go test -v -run TestYAML       # Single test function

# Run linters
make lint                      # Runs golangci-lint and govulncheck

# Format code
make format                    # Runs goimports on all files

# Initialize/update build environment
make builder-init              # Pulls latest builder rules
```

The project uses golangci-lint with a custom configuration (`.golangci.yaml`) that enables 50+ linters including experimental, performance, and style checkers.

## Architecture Overview

### Core Concepts

**Table** is a row-based data structure that:
1. Accepts structs via `Write()` and accumulates them as rows
2. Outputs to text (column-aligned), JSON, or YAML formats via `Flush*()`
3. Automatically flushes and starts a new table when struct types change
4. Applies ANSI colors/styles only when output is a terminal

**Three-pass text rendering:**
1. First pass (`FlushText`): Determine column widths and cell contents without formatting
2. Second pass (`flush`): Apply ANSI styles and handle repeating values
3. Output: Column-aligned text with proper spacing accounting for ANSI escape sequences

### Key Components

- **table.go**: Main Table type with Write/Annotate/Flush methods
- **text.go**: Text rendering with column alignment and ANSI styling
- **json.go/yaml.go**: Alternative output formats
- **tags.go**: Struct tag parsing (copied from stdlib json package)
- **sgr/**: Subpackage for ANSI Select Graphic Rendition escape sequences
  - **sgr.go**: Core SGR parameters and Wrap/Wrapped types
  - **wrap.go**: Text wrapping with ANSI codes, accounting for escape sequences in width calculations

### Struct Tag Convention

The `table` struct tag controls field behavior:

```go
type Example struct {
    Field1 string `table:"CUSTOM_LABEL"`           // Custom column header
    Field2 string `table:"LABEL,omitempty"`        // Hide column if all values are zero
    Field3 string `table:"-"`                      // Skip field entirely
}
```

Without tags, CamelCase field names are converted to UPPERCASE_SNAKE_CASE headers (e.g., `MyField` → `MY_FIELD`).

### Color Customization via Wrapper Interface

Types can implement the `wrapper` interface to apply custom ANSI styling:

```go
type wrapper interface {
    Wrap() sgr.Wrapped
}
```

Example from the codebase:
```go
type rank int

func (r rank) Wrap() sgr.Wrapped {
    return sgr.Wrap(color.Green, r)
}
```

When the Table encounters a value implementing `wrapper`, it calls `Wrap()` and uses the returned `sgr.Wrapped` for styled output. This only applies to text output mode.

### Terminal Detection and Color Handling

- Colors are disabled automatically if output is not a terminal (`term.IsTerminal`)
- The `sgr` package respects the `NO_COLOR` environment variable
- When colors are disabled, `sgr.Wrapped` types fall back to plain text via their `Text` field

## Conventions

### Import Paths
Always use `endobit.io/table` (not `github.com/endobit/table`) for imports.

### Testing Patterns
- Tests use real struct types (e.g., `host`) with varied field types
- Test functions often omit `*testing.T` parameter when not needed: `func TestYAML(_ *testing.T)`
- Coverage is tracked with `go test -coverprofile=coverage.out`

### Code Organization
- Each output format (text, JSON, YAML) has its own file
- Reflection-heavy code is isolated to processing functions
- ANSI/SGR logic is in a separate subpackage for reusability

### camelToUpperSnake Function
The default field-to-label conversion supports Unicode and handles acronyms intelligently:
- `CamelCase` → `CAMEL_CASE`
- `URLValue` → `URL_VALUE`
- `MyHTTPServer2` → `MY_HTTP_SERVER2`
- Supports non-ASCII: `ÖffentlicheVerkehrsmittel` → `ÖFFENTLICHE_VERKEHRSMITTEL`

### Annotations
Annotations are text strings inserted between table rows for comments/context. They:
- Only appear in text output (ignored in JSON/YAML)
- Are styled with the `Annotation` color scheme (default: Italic)
- Are added via `Annotate(string)` and tracked by row index
