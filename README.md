# SteuerGo

![License](https://img.shields.io/badge/license-MIT-blue)
![Go Version](https://img.shields.io/badge/go-%3E%3D1.18-blue)

A beautiful terminal-based German tax calculator using the official BMF (Bundesministerium der Finanzen) API.

## Overview

SteuerGo is a terminal application written in Go that helps users calculate their German income tax obligations. It uses the official BMF (Federal Ministry of Finance) API to ensure accurate tax calculations based on current German tax laws.

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
- 📝 Detailed tax information on demand
- 📅 Support for recent tax years

## Installation

### One-line installation

```bash
curl -sSL https://raw.githubusercontent.com/kyco/steuer-go/main/install.sh | bash
```

This will download the latest release binary for your platform and install it to `/usr/local/bin/steuergo`.

### From source

```bash
# Clone the repository
git clone https://github.com/kyco/steuer-go.git
cd steuergo

# Build the application
go build -o steuergo cmd/tax-calculator/main.go

# Run the application
./steuergo
```

### Using Go

```bash
go install github.com/kyco/steuer-go@latest
```

## Usage

Simply run the application and follow the on-screen instructions:

1. Select your tax class using the arrow keys
2. Enter your annual income
3. Confirm the tax year (default is the current year)
4. Press Tab to navigate between fields
5. Press Enter on the Calculate button to see your results

In the results screen:
- Press 'd' to toggle detailed tax information
- Press 'b' or 'Esc' to return to the input form
- Use arrow keys to scroll through results if needed

## Screenshots

[Coming Soon]

## How it works

SteuerGo connects to the official BMF API to calculate taxes based on the provided income and tax class. The API returns detailed tax information which is then formatted and displayed in a user-friendly way.

The application is built using:
- [Bubble Tea](https://github.com/charmbracelet/bubbletea): A powerful TUI framework
- [Lip Gloss](https://github.com/charmbracelet/lipgloss): For terminal styling
- Go's XML package: For parsing the BMF API responses

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## Acknowledgments

- The German Federal Ministry of Finance for providing the tax calculation API
- The Charm team for their excellent TUI libraries
