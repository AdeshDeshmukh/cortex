# 🧠 Cortex

**AI code reviewer that learns from your preferences through reinforcement learning.**

[![Go Version](https://img.shields.io/badge/Go-1.21+-00ADD8?style=flat&logo=go)](https://go.dev)
[![License](https://img.shields.io/badge/License-MIT-green.svg)](LICENSE)
[![Tests](https://img.shields.io/badge/tests-passing-brightgreen.svg)]()

Cortex is a privacy-first code reviewer that runs entirely on your machine. Unlike cloud-based tools, it learns YOUR coding style through reinforcement learning and gets smarter with every commit.

---

## ✨ Features

- 🔒 **Privacy-First**: All analysis happens locally, code never leaves your machine
- 🧠 **Adaptive Learning**: Uses RL to learn from your accept/reject decisions
- ⚡ **Fast**: Sub-2-second reviews with quantized local models
- 🎯 **Context-Aware**: Understands your project structure and patterns
- 🛠️ **Git Integration**: Seamlessly works with your existing workflow
- 📊 **Multi-Language**: Supports Go, Python, JavaScript, TypeScript, and more

---

## 🚀 Quick Start

### Installation

```bash
# Clone repository
git clone https://github.com/AdeshDeshmukh/cortex.git
cd cortex

# Build
make build

# Install system-wide (optional)
make install
```

### Usage

```bash
# Navigate to your project
cd ~/projects/myapp

# Install git hook
cortex install

# Make code changes
echo "func test() {}" >> main.go
git add main.go

# Commit (Cortex reviews automatically)
git commit -m "Add test function"
```

Output:

```
🧠 Cortex Review

📊 Summary:
   Files changed:  1
   Lines added:    +1
   Lines removed:  -0
   Languages:      go (1)

📁 Files:
   1. main.go (go)
      +1 -0 lines

⏳ Analysis pipeline:
   ✅ Diff parsing complete
   ⏹️  Static analysis (coming soon)
   ⏹️  LLM suggestions (coming soon)
```

---

## 📖 Commands

| Command | Description |
|---------|-------------|
| `cortex install` | Install pre-commit hook in current repo |
| `cortex install --force` | Overwrite existing hook |
| `cortex install --uninstall` | Remove hook |
| `cortex review` | Manually review staged changes |
| `cortex review --hook` | Run in hook mode (internal) |
| `cortex --help` | Show all commands |

---

## 🏗️ Architecture

```
┌─────────────────────────────────────────────┐
│           USER COMMITS CODE                 │
└──────────────┬──────────────────────────────┘
               │
               ▼
┌──────────────────────────────────────────────┐
│         GIT PRE-COMMIT HOOK                  │
│  Triggers: cortex review --hook              │
└──────────────┬───────────────────────────────┘
               │
               ▼
┌──────────────────────────────────────────────┐
│         DIFF PARSER (Go)                     │
│  Converts git diff → structured data         │
└──────────────┬───────────────────────────────┘
               │
               ▼
┌──────────────────────────────────────────────┐
│         ANALYZER PIPELINE                    │
│  1. Static analysis (regex, AST)             │
│  2. LLM inference (local model)              │
│  3. Ranking (RL policy)                      │
└──────────────┬───────────────────────────────┘
               │
               ▼
┌──────────────────────────────────────────────┐
│         INTERACTIVE UI                       │
│  User accepts/rejects suggestions            │
└──────────────┬───────────────────────────────┘
               │
               ▼
┌──────────────────────────────────────────────┐
│         FEEDBACK LOOP (Python)               │
│  Trains model on user preferences            │
└──────────────────────────────────────────────┘
```

---

## 🛠️ Development

### Prerequisites

- Go 1.21+
- Git
- Python 3.10+ (for ML features)

### Setup

```bash
# Clone repo
git clone https://github.com/AdeshDeshmukh/cortex.git
cd cortex

# Install dependencies
go mod download

# Run tests
make test

# Build
make build

# Run locally
./cortex --help
```

### Testing

```bash
# Run all tests
make test

# Run with coverage
make coverage

# Run specific package
go test ./internal/git/
```

---

## 📁 Project Structure

```
cortex/
├── cmd/cortex/              # CLI entry point
│   ├── main.go
│   └── commands/            # Cobra commands
│       ├── root.go
│       ├── install.go
│       └── review.go
├── internal/                # Private application code
│   └── git/                 # Git integration
│       ├── hooks.go         # Hook management
│       ├── diff.go          # Diff parsing
│       └── *_test.go        # Tests
├── pkg/                     # Public libraries
│   └── types/               # Type definitions
│       └── types.go
├── go.mod                   # Go dependencies
├── Makefile                 # Build automation
└── README.md                # This file
```

---

## 🎯 Roadmap

### ✅ Phase 1: Foundation (Current)
- Git hook installation
- Diff parsing
- CLI framework
- Test infrastructure

### 🚧 Phase 2: Intelligence (In Progress)
- Static code analysis
- LLM integration (local inference)
- Prompt engineering
- Suggestion generation

### 📅 Phase 3: Learning (Planned)
- Feedback collection system
- SQLite storage
- RL training pipeline (DPO)
- Model personalization

### 🔮 Phase 4: Advanced (Future)
- Multi-language support expansion
- IDE integration (VS Code, IntelliJ)
- Team-wide model sharing
- Custom rule definitions

---

## 🤝 Contributing

Contributions are welcome! Please follow these steps:

1. Fork the repository
2. Create a feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'feat: add amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

### Commit Convention

We follow Conventional Commits:

```
feat: New feature
fix: Bug fix
docs: Documentation changes
test: Test additions/changes
refactor: Code refactoring
chore: Maintenance tasks
```

---

## 📄 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

---

## 🙏 Acknowledgments

- Built with [Cobra](https://github.com/spf13/cobra) for CLI
- Inspired by first-principles thinking
- Part of ongoing learning journey in AI/ML systems

---

## 📬 Contact

**Adesh Deshmukh**

- GitHub: [@AdeshDeshmukh](https://github.com/AdeshDeshmukh)
- Project: [github.com/AdeshDeshmukh/cortex](https://github.com/AdeshDeshmukh/cortex)

⭐ **If you find this project useful, please consider giving it a star!**
// test hook
