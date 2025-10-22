# Developer Onboarding Guide

## Repository Layout

```
├─ conanfile.py           # Third-party libs
├─ cmake/
│   └─ profiles/           # Conan + CMake toolchain profiles
├─ config/
│   ├─ requirements.txt    # MkDocs + plugins
│   └─ *.md                # Documentation sources
├─ src/                    # C++ source tree
├─ tests/                  # Unit / integration tests
├─ Makefile                # High-level build driver
├─ README.md
```

## Documentation Workflow

### Prerequisites

Python packages from `config/requirements.txt` (already installed in the Prerequisites section).

### Build the Site
```
make config
```

- Generates a static site under `build/config`.
- Commit any updated Markdown files; the generated site can be committed optionally (CI can rebuild).

### Serve Locally (Live Reload)
```
make config-dev
```

- Starts a MkDocs server at [http://127.0.0.1:8000](http://127.0.0.1:8000/).
- The server watches `config/` and reloads on changes.

## Handy One-Liners (Copy-Paste)

### Clone & initial setup
```shell
git clone --recursive https://github.com/LetheanVPN/blockchain.git
cd blockchain
make configure
```

### Show Make targets
```shell
make help
```

### Start a feature branch
```shell
git checkout dev
git pull origin dev
git checkout -b dev-123456-my-feature
```

### Build a testnet binary
```shell
make gcc-linux-amd64 TESTNET=1
```

### Serve config locally
```shell
make config-dev
```

### Run all checks before opening a PR
```shell
make clean && make gcc-linux-amd64 TESTNET=1
make test
make config
```