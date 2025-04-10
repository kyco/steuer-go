# SteuerGo

![License](https://img.shields.io/badge/license-MIT-blue)
![Go Version](https://img.shields.io/badge/go-%3E%3D1.18-blue)
![GitHub release (latest by date)](https://img.shields.io/github/v/release/kyco/steuer-go)

A beautiful terminal-based German tax calculator using the official BMF (Bundesministerium der Finanzen) API, with support for offline calculation.

## Overview

SteuerGo is a terminal application written in Go that helps users calculate their German income tax obligations. It primarily uses the official BMF (Federal Ministry of Finance) API to ensure accurate tax calculations based on current German tax laws, but also provides a local calculation option that implements the official tax formulas for offline use.

With its clean, minimalist terminal interface (powered by Bubble Tea), SteuerGo makes it easy to:

- Calculate income tax based on different tax classes (Steuerklasse 1-6)
- View detailed tax breakdowns including income tax and solidarity tax
- See monthly and annual calculations side-by-side
- Visualize the proportion of taxes to net income

## Features

- 🖥️ Clean, intuitive terminal UI
- 🔢 Support for all German tax classes (1-6)
- 📊 Visual breakdown of tax calculations
- 🔄 Real-time calculations via the official BMF API
- 🧮 Offline calculation mode with local implementation of tax formula
- 📝 Detailed tax information on demand
- 📅 Support for recent tax years
- 📊 Comparative tax analysis across income levels

## Installation

### One-line installation

```bash
curl -sSL https://raw.githubusercontent.com/kyco/steuer-go/main/install.sh | bash
```

This will download the latest release binary for your platform and install it to `/usr/local/bin/steuergo`. The script supports macOS (Intel and ARM), Linux, and Windows.

### From source

```bash
# Clone the repository
git clone https://github.com/kyco/steuer-go.git
cd steuer-go

# Build the application
go build -o steuergo cmd/tax-calculator/main.go

# Run the application
./steuergo
```

### Using Go

```bash
go install github.com/kyco/steuer-go/cmd/tax-calculator@latest
```

After installation, the application will be available as `tax-calculator`. You can rename it to `steuergo` if you prefer:

```bash
mv $(which tax-calculator) $(dirname $(which tax-calculator))/steuergo
```

## Usage

After installation, run the application by typing `steuergo` in your terminal:

```bash
steuergo
```

Follow the on-screen instructions:

1. Select your tax class using the arrow keys
2. Enter your annual income
3. Confirm the tax year (default is the current year)
4. Press Tab to navigate between fields
5. Press Enter on the Calculate button to see your results

In the results screen:
- Press 'd' to toggle detailed tax information
- Press 'c' to compare tax rates across different income levels
- Press 'l' to toggle between online (API) and offline (local) calculation modes
- Press 'b' or 'Esc' to return to the input form
- Use arrow keys to scroll through results if needed

## Screenshots

[Coming Soon]

## How it works

SteuerGo primarily connects to the official BMF API to calculate taxes based on the provided income and tax class. The API returns detailed tax information which is then formatted and displayed in a user-friendly way.

For offline use or when the API is unavailable, SteuerGo can also perform calculations locally by implementing the German tax formula according to the official algorithm published by the BMF. This is based on the XML pseudo-code (PAP - Programmablaufplan) provided by the German tax authorities.

The application is built using:
- [Bubble Tea](https://github.com/charmbracelet/bubbletea): A powerful TUI framework
- [Lip Gloss](https://github.com/charmbracelet/lipgloss): For terminal styling
- Go's XML package: For parsing the BMF API responses and tax algorithm definition
- Custom implementation of the official German tax calculation formula

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## Acknowledgments

- The German Federal Ministry of Finance for providing the tax calculation API
- The Charm team for their excellent TUI libraries
