# 🧰 ghpp — GitHub Pretty Printer

`ghpp` is a lightweight CLI tool written in Go that fetches and pretty-prints metadata for any public GitHub repository. It supports colorful and emoji-rich output, configurable fields, and human-readable timestamps.

---

## ✨ Features

- 📦 Displays key GitHub repo details (name, URL, stars, etc.)
- 📝 Supports `--include` and `--exclude` fields
- ⏳ Shows `created_at` as relative time (`--since`, `--short`)
- 🌈 Colorful terminal output with emoji icons
- 🔌 Easily extensible and configurable

---

## 🚀 Installation

### 🔧 Prerequisites

- Go 1.20 or later
- Internet connection (for GitHub API)

### 🛠️ Build & Install

```bash
git clone https://github.com/kirancnayak/ghpp.git
cd ghpp

# To build a binary
make build

# To install globally (adds it to $GOBIN)
make install
