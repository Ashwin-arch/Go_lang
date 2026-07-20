# Go_lang

A repository for Go language projects and applications.

---

## 🧮 Calculator (GUI)

A cross-platform desktop GUI calculator built using Go and the [Fyne v2](https://fyne.io/) GUI toolkit.

### Features
- Standard arithmetic operations: addition (`+`), subtraction (`-`), multiplication (`×`), division (`÷`), percentage (`%`)
- Sign toggle (`±`), clear (`C`), and backspace (`⌫`)
- Modern dark-themed user interface

---

## 🚀 How to Run

### 1. Prerequisites
Ensure you have [Go](https://go.dev/doc/install) (v1.22 or higher) installed on your system.

#### Linux System Dependencies (for Fyne GUI)
If running on Linux, Fyne requires standard C compiler graphics development packages:
- **Ubuntu/Debian**:
  ```bash
  sudo apt update
  sudo apt install golang gcc libgl1-mesa-dev xorg-dev
  ```
- **Fedora**:
  ```bash
  sudo dnf install golang gcc libX11-devel libXcursor-devel libXrandr-devel libXinerama-devel mesa-libGL-devel libXi-devel libXxf86vm-devel
  ```
- **Arch Linux**:
  ```bash
  sudo pacman -S go gcc libxcursor libxrandr libxinerama libxi mesa
  ```

---

### 2. Running the Application

1. **Clone the Repository** (if not already local):
   ```bash
   git clone https://github.com/Ashwin-arch/Go_lang.git
   cd Go_lang/calculator
   ```

2. **Download Dependencies**:
   ```bash
   go mod download
   ```

3. **Run Directly**:
   ```bash
   go run main.go
   ```
   *or*
   ```bash
   go run .
   ```

---

### 3. Building an Executable

To compile a standalone binary executable:

```bash
go build -o calculator main.go
```

Then launch the executable:
- **Linux / macOS**: `./calculator`
- **Windows**: `calculator.exe`
