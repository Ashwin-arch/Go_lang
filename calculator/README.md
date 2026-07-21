# 🧮 Go GUI Calculator

A cross-platform desktop GUI calculator application written in Go using the [Fyne v2](https://fyne.io/) GUI toolkit.

---

## 📋 Table of Contents
1. [Prerequisites](#-prerequisites)
2. [How to Run & Build](#-how-to-run--build)
3. [Features](#-features)
4. [Detailed Code Explanations & Architecture (`main.go`)](#-detailed-code-explanations--architecture-maingo)

---

## ⚙️ Prerequisites

- **Go 1.22+** installed on your system (`go version`).
- **Graphics Dependencies** (Fyne requires OpenGL driver & C compiler):
  - **Linux**: `sudo apt install gcc libgl1-mesa-dev xorg-dev`
  - **macOS**: Xcode Command Line Tools (`xcode-select --install`)
  - **Windows**: MinGW-w64 compiler or gcc in PATH

---

## 🚀 How to Run & Build

### 1. Run Directly from Source
```bash
cd calculator
go run .
```

### 2. Build Standalone Executable Binary
```bash
cd calculator
go build -o calculator main.go
./calculator
```

---

## ✨ Features

- Standard arithmetic operations (`+`, `-`, `×`, `÷`, `%`)
- Sign toggle (`±`), clear (`C`), and backspace (`⌫`)
- Divide-by-zero protection returning `Error: Div by 0`
- Modern dark-themed user interface with color-coded button importance levels

---

## 📂 Detailed Code Explanations & Architecture (`main.go`)

### 1. Calculator State (`calculator` struct)
```go
type calculator struct {
    display *widget.Entry
    op      string
    val1    float64
    newNum  bool
}
```
- **`display`**: Pointer to Fyne's `widget.Entry`, displaying user numbers and calculated results.
- **`op`**: Holds active arithmetic operation string (`"+"`, `"-"`, `"×"`, `"÷"`, `"%"`).
- **`val1`**: Holds the initial operand value parsed when an operator is pressed.
- **`newNum`**: Flag determining whether entering a digit overwrites the current display text or appends to it.

### 2. Core Calculator Operations
- **`input(digit string)`**: Handles numeric input and decimal points, ensuring only a single decimal point can be entered.
- **`setOp(op string)`**: Parses display text into `val1`, triggers calculation if another operator was already active (e.g. `10 + 5 - 3`), sets `op`, and marks `newNum = true`.
- **`compute()`**: Parses `val2`, executes operation in a `switch` block, handles divide-by-zero (`val2 == 0`), and formats floating point output cleanly using `strconv.FormatFloat(res, 'f', -1, 64)`.
- **`toggleSign()`, `backspace()`, `clear()`**: Helper methods for sign flipping (`±`), single-digit deletion (`⌫`), and state reset (`C`).

### 3. Fyne GUI Layout & Styling
- **`app.NewWithID(...)` & Dark Theme**: Sets theme to `theme.DarkTheme()` for a sleek UI look.
- **Button Factory (`makeButton`)**: Custom helper initializing buttons with specific Fyne importance levels:
  - `widget.HighImportance`: Primary action buttons (Operators `÷`, `×`, `-`, `+`, `=`).
  - `widget.MediumImportance`: Numeric buttons (`0`-`9`, `.`).
  - `widget.LowImportance`: Auxiliary control buttons (`C`, `⌫`, `±`, `%`).
- **Grid Layout (`container.NewGridWithColumns(4, ...)`)**: Renders a 4-column button layout.
- **Border Layout (`container.NewBorder(...)`)**: Docking container placing the display bar at the top and filling remaining window space with the button grid.
