# Studex Portal CLI

Studex Portal CLI is the command-line interface version of **[Studex Portal](https://studex.itshivam.in)**, Social Media Platform. 

It allows you to bring the complete functionality of the Studex Portal social networking experience directly into your terminal. Designed to be fast, extensible, and flexible, the CLI is published entirely across all familiar software package environments!

## Installation

Studex CLI provides deep multi-language and OS-native packaging layers. You can seamlessly install it via your favorite package manager.

### JavaScript / Node.js (NPM)
```sh
npm install -g studex-cli

# Or use it directly without installing globally:
npx studex-cli
```

### Python (PIP)
```sh
pip install studex-cli
studex-cli --help
```

### macOS / Linux (Homebrew)
```sh
brew tap itshivams/homebrew-studex-cli
brew install studex-cli
```

### Debian / Ubuntu (APT)
Download the latest `.deb` package from our [Releases](https://github.com/itshivams/Studex-CLI/releases) page:
```sh
sudo apt install ./studex-cli_*_linux_amd64.deb
```

### Fedora / RHEL / CentOS (YUM / DNF)
Download the latest `.rpm` package from [Releases](https://github.com/itshivams/Studex-CLI/releases):
```sh
sudo dnf install ./studex-cli_*_linux_amd64.rpm
```

### Windows (Winget & Chocolatey)

**Winget**:
```powershell
winget install Studex.StudexCLI
```

**Chocolatey**:
```powershell
choco install studex-cli
```

*(Note: Check the GitHub Releases page for direct binary downloads if you prefer no package managers.)*

---

## Setup & Local Development

Want to test it out locally before creating a PR?

1. Ensure you have [Go (v1.21+)](https://go.dev/doc/install) installed.
2. Clone the repository:
   ```sh
   git clone https://github.com/itshivams/Studex-CLI.git
   cd Studex-CLI
   ```
3. Run the CLI directly:
   ```sh
   go run main.go
   ```
4. Build the binary manually:
   ```sh
   go build -o studex-cli main.go
   ```

---

## Contributing (Open Source)

We highly encourage open source contributions! We want Studex CLI to adapt to multiple platforms, continuously improve features, and fix issues dynamically. 

### Branching Strategy

Please follow this branching convention when creating your PRs to ensure an organized workflow:

- `dev/<username>` : For overarching workspace modifications by a developer.
- `feature/*` : For engineering new features (e.g., `feature/user-auth`).
- `fix/*` : For standard bug fixes (e.g., `fix/search-timeout`).
- `hotfix/*` : For urgent production fixes (e.g., `hotfix/crash-on-login`).

### Submission Workflow

1. Fork the repository.
2. Create your branch adapting the above strategy (`git checkout -b feature/amazing-feature`).
3. Commit your changes (`git commit -m 'Add: xyz feature'`).
4. Push to the branch (`git push origin feature/amazing-feature`).
5. Open a **Pull Request**.

**CI Status Note:** Our `.github/workflows/ci.yml` strictly guards the `main` branch. Every Pull Request will automatically be tested against all supported OS environments (Windows, macOS, Linux) and will thoroughly validate the integrity of cross-language packages (`npm`, `pip`, `apt`, `winget`). Ensure tests pass before pushing!


## License

This project is licensed under the [MIT License](LICENSE).
