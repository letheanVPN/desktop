---
title: Core.Help
---

# Overview

Core is an opinionated framework for building Go desktop apps with Wails, providing a small set of focused modules you can mix into your app. It ships with sensible defaults and a demo app that doubles as in‑app help.

- Site: [https://dappco.re](https://dappco.re)
- Repo: [https://github.com/Snider/Core](https://github.com/Snider/Core)

## Modules

- Core — framework bootstrap and service container
- Core.Config — app and UI state persistence
- Core.Crypt — keys, encrypt/decrypt, sign/verify
- Core.Display — windows, tray, window state
- Core.Help — in‑app help and deep‑links
- Core.IO — local/remote filesystem helpers
- Core.Workspace — projects and paths

## Quick start
```go
package main

import (
    "github.com/wailsapp/wails/v3/pkg/application"
    core "github.com/Snider/Core"
)

func main() {
    app := core.New(
        core.WithServiceLock(),
    )
    wailsApp := application.NewWithOptions(&application.Options{
        Bind: []interface{}{app},
    })
    wailsApp.Run()
}
```

## Services
```go
package demo

import (
    core "github.com/Snider/Core"
)

// Register your service
func Register(c *core.Core) error {
    return c.RegisterService("demo", &Demo{core: c})
}
```

## Display example
```go
package display

import (
    "context"
    "github.com/wailsapp/wails/v3/pkg/application"
)

// Open a window on startup
func (d *API) ServiceStartup(ctx context.Context, _ application.ServiceOptions) error {
    d.OpenWindow(
        OptName("main"),
        OptHeight(900),
        OptWidth(1280),
        OptURL("/"),
        OptTitle("Core"),
    )
    return nil
}
```

See the left nav for detailed pages on each module.
