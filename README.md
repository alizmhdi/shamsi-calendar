# Scal - Jalali (Shamsi) Calendar CLI

A command-line tool to display Jalali (Shamsi) calendar, similar to the Unix `cal` command.

## Installation

### From Source

```bash
git clone https://github.com/alizmhdi/shamsi-calendar.git
cd shamsi-calendar

go build -o scal .

sudo mv scal /usr/local/bin/
```

### Using Go Install

```bash
go install github.com/alizmhdi/shamsi-calendar/cmd/scal@latest
```

## Usage

### Basic Commands

```bash
# Display current month
scal

# Display specific month
scal -m 4

# Display specific year (full year)
scal -y 1404

# Display specific month and year
scal -y 1404 -m 4

# Display three months (previous, current, next)
scal -3

# Display full year for current year
scal -Y
```

### Command Line Options

| Flag | Short | Description | Example |
|------|-------|-------------|---------|
| `--year` | `-y` | Year to display (default: current year) | `scal -y 1404` |
| `--month` | `-m` | Month to display (1-12, default: current month) | `scal -m 4` |
| `--three` | `-3` | Display three months spanning the date | `scal -3` |
| `--full-year` | `-Y` | Display entire year | `scal -Y` |
