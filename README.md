# webpage

A Go library and CLI tool for capturing webpages as PDFs or images.

## Supported Features

- PDFs
- Images (PNG/JPEG)

## Requirements

- Chrome/Chromium browser installed on the system

## Installation

### Library

```bash
go get github.com/joeychilson/webpage
```

### CLI Tool

```bash
go install github.com/joeychilson/webpage/cmd/webpage@latest
```

## Library Usage

### Basic PDF Generation

```go
package main

import (
    "context"
    "os"
    "github.com/joeychilson/webpage"
)

func main() {
    page := webpage.New("https://example.com")
    
    pdf, err := page.PDF(context.Background())
    if err != nil {
        panic(err)
    }
    
    os.WriteFile("output.pdf", pdf, 0644)
}
```

### Customized PDF Generation

```go
page := webpage.New("https://example.com",
    webpage.WithTimeout(60*time.Second),
    webpage.WithUserAgent("Custom User Agent"),
)

pdf, err := page.PDF(context.Background(),
    webpage.WithLandscape(true),
    webpage.WithBackground(true),
    webpage.WithScale(1.2),
    webpage.WithPaperWidth(11.0),
    webpage.WithPaperHeight(17.0),
    webpage.WithMarginTop(1.0),
    webpage.WithPageRanges("1-5"),
)
```

### Screenshot Capture

```go
page := webpage.New("https://example.com")

screenshot, err := page.Screenshot(context.Background(),
    webpage.WithFormat("jpeg"),
    webpage.WithQuality(90),
)
```

## CLI Usage

The `webpage` CLI tool provides a command-line interface to the webpage capture functionality.

### Basic Commands

```bash
# Generate PDF
webpage pdf https://example.com -o output.pdf

# Capture screenshot
webpage screenshot https://example.com -o output.png
```

### PDF Options

```bash
webpage pdf https://example.com -o output.pdf \
    --landscape \
    --background \
    --scale 1.2 \
    --width 11 \
    --height 17 \
    --margin-top 1 \
    --margin-bottom 1 \
    --margin-left 1 \
    --margin-right 1 \
    --pages "1-5"
```

### Screenshot Options

```bash
webpage screenshot https://example.com -o output.jpg \
    --format jpeg \
    --quality 90
```

### Common Options

```bash
# Set timeout
webpage pdf https://example.com -o output.pdf --timeout 60s

# Set user agent
webpage pdf https://example.com -o output.pdf --user-agent "Custom User Agent"
```

## Configuration Options

### Browser Options

| Option | CLI Flag | Description | Default |
|--------|----------|-------------|---------|
| Timeout | `--timeout, -t` | Operation timeout | 30s |
| User Agent | `--user-agent, -u` | Browser user agent | Chrome default |

### PDF Options

| Option | CLI Flag | Description | Default |
|--------|----------|-------------|---------|
| Landscape | `--landscape, -l` | Landscape orientation | false |
| Background | `--background, -b` | Include background graphics | false |
| Scale | `--scale, -s` | Scale factor | 1.0 |
| Paper Width | `--width` | Paper width (inches) | 8.5 |
| Paper Height | `--height` | Paper height (inches) | 11.0 |
| Margin Top | `--margin-top` | Top margin (inches) | 0.4 |
| Margin Bottom | `--margin-bottom` | Bottom margin (inches) | 0.4 |
| Margin Left | `--margin-left` | Left margin (inches) | 0.4 |
| Margin Right | `--margin-right` | Right margin (inches) | 0.4 |
| Page Ranges | `--pages` | Page ranges to print | "" (all) |

### Screenshot Options

| Option | CLI Flag | Description | Default |
|--------|----------|-------------|---------|
| Format | `--format, -f` | Image format (png/jpeg) | png |
| Quality | `--quality, -q` | JPEG quality (0-100) | 100 |

