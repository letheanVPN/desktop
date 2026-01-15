# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

Lethean Desktop is a cross-platform desktop application built with **Wails v3** (Go backend + Angular frontend). It provides decentralized VPN and blockchain functionality for the Lethean network.

## Build & Development Commands

```bash
# Development mode (hot-reload for both Go and Angular)
wails3 dev -config ./build/config.yml -port 9245
# or use Task
task dev

# Build production binary
task build
# or
wails3 build

# Package for distribution (.app on macOS)
task package

# Build universal macOS binary (arm64 + amd64)
task darwin:build:universal

# Generate TypeScript bindings from Go services
wails3 generate bindings -ts -i

# Run Go tests
go test ./services/...
# or
make test

# Run a single test file/package
go test ./services/core/config/...
go test -run TestNewService ./services/core/config/...

# Frontend commands (run from ./frontend/)
npm install          # Install dependencies
npm run dev          # Dev server (port 4200)
npm run build        # Production build
npm run build:dev    # Development build
npm test             # Run Karma tests
```

## Architecture

### Backend (Go)

**Two-phase initialization pattern** in `main.go`:
1. `core.New(assets)` - Instantiates all services without dependencies
2. `core.Setup(app)` - Wires services together after Wails app is created

**Core services** (`services/core/`):
- `core.go` - Singleton service container with accessor functions (`core.Config()`, `core.I18n()`, etc.)
- `config/` - Application configuration and XDG-compliant paths
- `display/` - Window management, system tray, menus
- `i18n/` - Internationalization using go-i18n

**Application services** (`services/`):
- `blockchain/` - Network interface with pluggable implementations (Lethean network)
- `mining/` - Mining service with HTTP API (gin router on :8080), supports xmrig/ttminer
- `crypt/` - OpenPGP encryption, hashing utilities
- `filesystem/` - Abstraction layer for local/SFTP/WebDAV storage
- `workspace/` - User workspace management

**Service registration**: Services implement Wails service lifecycle (`ServiceStartup`, etc.) and are registered in `main.go` via `application.NewService()`.

### Frontend (Angular)

**Stack**: Angular 20, TailwindCSS 4, FontAwesome/WebAwesome, Highcharts, Monaco Editor

**Key directories**:
- `src/app/` - Components and services
- `src/frame/` - Application frames (main UI, system tray)
- `bindings/` - Auto-generated TypeScript bindings from Go services

**Calling Go from Angular**: Import from `@lthn/*` paths which map to generated bindings:
```typescript
import { ShowEnvironmentDialog } from "@lthn/core/display/service"
import { IsFeatureEnabled } from "@lthn/core/config/service"
```

**Routes** (`app.routes.ts`):
- `/` - Main application frame with blockchain, mining, dev views
- `/setup/*` - Setup wizard components
- `/system-tray` - System tray popup
- `/editor/monaco` - Code editor

### Build System

**Taskfile.yml** hierarchy:
- Root `Taskfile.yml` - Main entry points (`build`, `dev`, `package`, `run`)
- `build/Taskfile.yml` - Common tasks (bindings, frontend build, icons)
- `build/{darwin,linux,windows}/Taskfile.yml` - Platform-specific builds

**Key build flags**:
- Production: `-tags production -trimpath -ldflags="-w -s"`
- Development: `-gcflags=all="-l"`

## Frontend/Backend Binding

Go services expose methods to the frontend via Wails bindings. The bindings are generated into `frontend/bindings/` and imported with `@lthn/` prefix. When adding new Go service methods, run `wails3 generate bindings -ts -i` to update TypeScript types.

## i18n

- Backend: JSON locale files in `services/core/i18n/locales/`
- Frontend: TranslationService wraps @ngx-translate, loads from `/assets/i18n/`
- Both use the same JSON files (copied via angular.json assets config)
