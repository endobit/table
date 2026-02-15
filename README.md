# table

[![GitHub Actions](https://img.shields.io/github/actions/workflow/status/endobit/table/test.yaml)](https://github.com/endobit/table/actions?query=workflow%3Atest)
[![Go Version](https://img.shields.io/github/go-mod/go-version/endobit/table)](https://img.shields.io/github/go-mod/go-version/endobit/table)
[![Go Report Card](https://goreportcard.com/badge/github.com/endobit/table)](https://goreportcard.com/report/endobit.io/table)
[![Codecov](https://codecov.io/gh/endobit/oui/branch/main/graph/badge.svg)](https://codecov.io/gh/endobit/table)
[![Go Reference](https://pkg.go.dev/badge/github.com/endobit/table.svg)](https://pkg.go.dev/endobit.io/table)

A Go package for rendering structs as column-aligned tables with optional ANSI color styling, or as JSON/YAML.

## Features

- **Text Output**: Column-aligned tables with automatic header generation
- **Multiple Formats**: Output as text, JSON, or YAML
- **ANSI Styling**: Automatic color/style support with terminal detection
- **Struct Tags**: Customize column headers and behavior with `table` tags
- **Annotations**: Insert comments between table rows
- **Custom Colors**: Apply custom ANSI styling via the `wrapper` interface
- **Smart Defaults**: CamelCase field names convert to UPPERCASE_SNAKE_CASE headers

## Installation

```bash
go get endobit.io/table
```

## Usage

### Basic Example

```go
package main

import (
    "endobit.io/table"
)

type server struct {
    Name   string
    Status string
    Port   int
}

func main() {
    t := table.New()
    
    t.Write(server{Name: "web-1", Status: "running", Port: 8080})
    t.Write(server{Name: "web-2", Status: "stopped", Port: 8081})
    
    _ = t.Flush()
}
```

Output:
```
NAME  STATUS  PORT
web-1 running 8080
web-2 stopped 8081
```

### Struct Tags

Use the `table` tag to customize column headers or hide fields:

```go
type host struct {
    Zone    string `table:"ZONE"`              // Custom header
    Cluster string `table:"CLUSTER"`
    Host    string `table:"HOST"`
    Rack    string `table:"RACK,omitempty"`    // Hide if all values are zero
    Rank    int    `table:"RANK"`
    Internal string `table:"-"`                // Skip this field
}
```

### Annotations

Insert comments or context between rows:

```go
t := table.New()
t.Write(host{Zone: "east", Cluster: "prod", Host: "compute-0"})
t.Annotate("maintenance window scheduled")
t.Write(host{Zone: "west", Cluster: "prod", Host: "compute-1"})
_ = t.Flush()
```

### JSON and YAML Output

```go
t := table.New()
t.Write(server{Name: "web-1", Status: "running", Port: 8080})

// Output as JSON
_ = t.FlushJSON()

// Output as YAML
_ = t.FlushYAML()
```

### Custom ANSI Colors

Implement the `wrapper` interface to apply custom styling:

```go
import "endobit.io/table/sgr/color"

type rank int

func (r rank) Wrap() sgr.Wrapped {
    return sgr.Wrap(color.Green, r)
}

type host struct {
    Name string
    Rank rank  // Will be rendered in green
}
```

### Configuration Options

```go
// Custom writer
t := table.New(table.WithWriter(myWriter))

// Custom colors
colors := &table.Colors{
    Header:  []sgr.Param{sgr.Bold, sgr.Underline},
    EvenRow: []sgr.Param{sgr.Faint},
}
t := table.New(table.WithColor(colors))

// Custom label function
t := table.New(table.WithLabelFunction(strings.ToLower))
```

## Color Handling

- Colors automatically disabled when output is not a terminal
- Respects `NO_COLOR` environment variable
- ANSI escape sequences properly handled in column width calculations

## Examples

See the [example](example/) directory for a complete working example.
