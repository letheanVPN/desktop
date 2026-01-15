# Electron Compatibility Layer for Wails

This module provides Electron-like APIs that map to Wails v3 runtime equivalents. It's designed to help developers familiar with Electron contribute to this Wails-based application.

## Quick Start

```typescript
import { ipcRenderer, shell, dialog, app, clipboard, BrowserWindow } from '@lib/electron-compat';

// IPC Communication - just like Electron!
const result = await ipcRenderer.invoke('blockchain:fetchBlockData', '12345');

// Open external links
await shell.openExternal('https://lethean.io');

// File dialogs
const files = await dialog.showOpenDialog({
  title: 'Select Wallet',
  filters: [{ name: 'Wallet', extensions: ['wallet', 'keys'] }]
});

// Clipboard
await clipboard.writeText(walletAddress);

// Window management
const win = BrowserWindow.getFocusedWindow();
win.setTitle('Lethean Desktop');
win.maximize();
```

## API Mapping Reference

### ipcRenderer

| Electron | Wails | Notes |
|----------|-------|-------|
| `ipcRenderer.send(channel, ...args)` | `Events.Emit()` | Fire-and-forget |
| `ipcRenderer.invoke(channel, ...args)` | `Call.ByName()` | Returns Promise |
| `ipcRenderer.on(channel, listener)` | `Events.On()` | Subscribe |
| `ipcRenderer.once(channel, listener)` | `Events.Once()` | One-time |
| `ipcRenderer.sendSync()` | ❌ | Not supported |

**Channel Naming Convention:**
```typescript
// Electron-style channels are auto-converted:
'blockchain:fetchBlockData' → 'blockchain.Service.FetchBlockData'
'config:get'                → 'config.Service.Get'

// Or use direct Wails binding paths:
'github.com/letheanVPN/desktop/services/blockchain.Service.FetchBlockData'
```

### shell

| Electron | Wails | Status |
|----------|-------|--------|
| `shell.openExternal(url)` | `Browser.OpenURL()` | ✅ Works |
| `shell.openPath(path)` | Go backend | ⚠️ Needs Go service |
| `shell.showItemInFolder(path)` | Go backend | ⚠️ Needs Go service |
| `shell.beep()` | Web Audio API | ✅ Works |

### dialog

| Electron | Wails | Status |
|----------|-------|--------|
| `dialog.showOpenDialog()` | `Dialogs.OpenFile()` | ✅ Works |
| `dialog.showSaveDialog()` | `Dialogs.SaveFile()` | ✅ Works |
| `dialog.showMessageBox()` | `Dialogs.Info/Warning/Error/Question()` | ✅ Simplified |
| `dialog.showErrorBox()` | `Dialogs.Error()` | ✅ Works |

### BrowserWindow

| Electron | Wails | Status |
|----------|-------|--------|
| `new BrowserWindow()` | Go `display.Service.OpenWindow()` | ⚠️ Via Go |
| `win.maximize/minimize()` | `Window.Maximise/Minimise()` | ✅ Works |
| `win.setTitle()` | `Window.SetTitle()` | ✅ Works |
| `win.setSize/Position()` | `Window.SetSize/Position()` | ✅ Works |
| `win.on(event)` | `Events.On()` | ✅ Works |
| Multi-window support | Go backend | ⚠️ Different model |

### clipboard

| Electron | Wails | Status |
|----------|-------|--------|
| `clipboard.readText()` | `navigator.clipboard` | ✅ Async |
| `clipboard.writeText()` | `navigator.clipboard` | ✅ Async |
| `clipboard.readImage()` | `navigator.clipboard` | ✅ Async |
| Sync methods | ❌ | Browser limitation |

### app

| Electron | Wails | Status |
|----------|-------|--------|
| `app.quit()` | `Application.Quit()` | ✅ Works |
| `app.getVersion()` | Hardcoded | ⚠️ Could bind |
| `app.getPath()` | Go backend | ⚠️ Needs Go service |
| `app.getLocale()` | `navigator.language` | ✅ Works |

## Key Differences from Electron

### 1. Process Model
- **Electron**: Main process + Renderer process(es)
- **Wails**: Go backend + Single frontend (WebView)

### 2. IPC Communication
- **Electron**: `ipcMain`/`ipcRenderer` with event-based messaging
- **Wails**: Direct Go method calls via bindings + Events for pub/sub

### 3. Window Management
- **Electron**: Create windows freely from main or renderer
- **Wails**: Windows created from Go backend, controlled via runtime

### 4. Native Features
- **Electron**: Node.js APIs available in renderer (with nodeIntegration)
- **Wails**: Native features exposed through Go bindings

## Backend Requirements

Some APIs require Go backend services. Create these in `services/core/`:

### shell/service.go (for shell.openPath, etc.)

```go
package shell

import (
    "os/exec"
    "runtime"
)

type Service struct{}

func NewService() *Service {
    return &Service{}
}

func (s *Service) OpenPath(path string) error {
    var cmd *exec.Cmd
    switch runtime.GOOS {
    case "darwin":
        cmd = exec.Command("open", path)
    case "linux":
        cmd = exec.Command("xdg-open", path)
    case "windows":
        cmd = exec.Command("cmd", "/c", "start", "", path)
    }
    return cmd.Start()
}

func (s *Service) ShowItemInFolder(path string) error {
    switch runtime.GOOS {
    case "darwin":
        return exec.Command("open", "-R", path).Start()
    case "linux":
        return exec.Command("xdg-open", filepath.Dir(path)).Start()
    case "windows":
        return exec.Command("explorer", "/select,", path).Start()
    }
    return nil
}
```

## Adding New Channel Mappings

Edit `ipc-renderer.ts` and add to `channelMappings`:

```typescript
const channelMappings: Record<string, string> = {
  // Add your mappings here
  'myService:myMethod': 'github.com/letheanVPN/desktop/services/mypackage.Service.MyMethod',
};
```

## Example: Migrating Electron Code

### Before (Electron)
```typescript
const { ipcRenderer, shell } = require('electron');

// Call main process
const data = await ipcRenderer.invoke('get-wallet-balance', walletId);

// Open link
shell.openExternal('https://explorer.lethean.io');

// File dialog
const { filePaths } = await ipcRenderer.invoke('show-open-dialog', {
  filters: [{ name: 'Wallet', extensions: ['wallet'] }]
});
```

### After (Wails + electron-compat)
```typescript
import { ipcRenderer, shell, dialog } from '@lib/electron-compat';

// Call Go service (same API!)
const data = await ipcRenderer.invoke('wallet:getBalance', walletId);

// Open link (identical!)
await shell.openExternal('https://explorer.lethean.io');

// File dialog (slightly different, dialog is in frontend)
const { filePaths } = await dialog.showOpenDialog({
  filters: [{ name: 'Wallet', extensions: ['wallet'] }]
});
```

## Contributing

When adding Electron API compatibility:

1. Check if Wails has a direct equivalent in `@wailsio/runtime`
2. If not, determine if it needs a Go backend service
3. Add proper TypeScript types
4. Document any behavioral differences
5. Add to this README

## License

EUPL-1.2 (same as parent project)
