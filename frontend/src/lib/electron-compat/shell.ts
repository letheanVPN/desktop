/**
 * Electron shell Compatibility Layer
 *
 * Maps Electron's shell API to Wails Browser/runtime equivalents.
 *
 * Electron Concept -> Wails Equivalent:
 * - shell.openExternal()  -> Browser.OpenURL()
 * - shell.openPath()      -> Runtime call to Go's os/exec
 * - shell.showItemInFolder() -> Runtime call to Go's file manager
 *
 * @example
 * import { shell } from '@lib/electron-compat';
 *
 * shell.openExternal('https://lethean.io');
 * shell.showItemInFolder('/path/to/file.txt');
 */

import { Browser, Call } from '@wailsio/runtime';

export const shell = {
  /**
   * Open a URL in the user's default browser.
   *
   * @param url - The URL to open
   * @param _options - Electron options (ignored in Wails)
   * @returns Promise that resolves when the URL is opened
   *
   * @example
   * await shell.openExternal('https://github.com/letheanVPN/desktop');
   * await shell.openExternal('mailto:support@lethean.io');
   */
  async openExternal(url: string, _options?: { activate?: boolean }): Promise<void> {
    Browser.OpenURL(url);
  },

  /**
   * Open a file or folder with the system's default application.
   *
   * NOTE: This requires a Go backend method to be implemented.
   * See the comment below for the Go implementation.
   *
   * @param path - The path to open
   * @returns Promise resolving to an error string (empty if success)
   *
   * @example
   * const error = await shell.openPath('/Users/me/Documents/file.pdf');
   * if (error) console.error('Failed to open:', error);
   */
  async openPath(path: string): Promise<string> {
    try {
      // This needs a Go backend method - see shell_backend.go below
      await Call.ByName('github.com/letheanVPN/desktop/services/core/shell.Service.OpenPath', path);
      return '';
    } catch (error) {
      return String(error);
    }
  },

  /**
   * Show a file in its parent folder with the file selected.
   *
   * NOTE: This requires a Go backend method to be implemented.
   * See the comment below for the Go implementation.
   *
   * @param fullPath - The full path to the file
   *
   * @example
   * shell.showItemInFolder('/Users/me/Downloads/blockchain.dat');
   */
  async showItemInFolder(fullPath: string): Promise<void> {
    try {
      // This needs a Go backend method - see shell_backend.go below
      await Call.ByName('github.com/letheanVPN/desktop/services/core/shell.Service.ShowItemInFolder', fullPath);
    } catch (error) {
      console.error('[electron-compat] showItemInFolder failed:', error);
    }
  },

  /**
   * Move a file to the system trash/recycle bin.
   *
   * NOTE: This requires a Go backend method to be implemented.
   *
   * @param fullPath - The full path to the file
   * @returns Promise resolving to void
   *
   * @example
   * await shell.trashItem('/Users/me/old-file.txt');
   */
  async trashItem(fullPath: string): Promise<void> {
    try {
      await Call.ByName('github.com/letheanVPN/desktop/services/core/shell.Service.TrashItem', fullPath);
    } catch (error) {
      throw new Error(`Failed to trash item: ${error}`);
    }
  },

  /**
   * Play the system beep sound.
   */
  beep(): void {
    // Use the Web Audio API as a fallback
    try {
      const audioContext = new (window.AudioContext || (window as any).webkitAudioContext)();
      const oscillator = audioContext.createOscillator();
      oscillator.type = 'sine';
      oscillator.frequency.value = 800;
      oscillator.connect(audioContext.destination);
      oscillator.start();
      oscillator.stop(audioContext.currentTime + 0.1);
    } catch {
      console.log('\u0007'); // ASCII bell character fallback
    }
  },

  /**
   * Read a shortcut file (Windows .lnk files).
   * Not applicable on macOS/Linux.
   *
   * @deprecated Platform-specific, not implemented in Wails
   */
  readShortcutLink(_shortcutPath: string): { target: string } {
    console.warn('[electron-compat] readShortcutLink is Windows-only and not implemented');
    return { target: '' };
  },

  /**
   * Write a shortcut file (Windows .lnk files).
   * Not applicable on macOS/Linux.
   *
   * @deprecated Platform-specific, not implemented in Wails
   */
  writeShortcutLink(_shortcutPath: string, _options: unknown): boolean {
    console.warn('[electron-compat] writeShortcutLink is Windows-only and not implemented');
    return false;
  },
};

/*
 * =============================================================================
 * GO BACKEND IMPLEMENTATION REQUIRED
 * =============================================================================
 *
 * Create this file at: services/core/shell/service.go
 *
 * ```go
 * package shell
 *
 * import (
 *     "os/exec"
 *     "runtime"
 * )
 *
 * type Service struct{}
 *
 * func NewService() *Service {
 *     return &Service{}
 * }
 *
 * // OpenPath opens a file or folder with the system default application.
 * func (s *Service) OpenPath(path string) error {
 *     var cmd *exec.Cmd
 *     switch runtime.GOOS {
 *     case "darwin":
 *         cmd = exec.Command("open", path)
 *     case "linux":
 *         cmd = exec.Command("xdg-open", path)
 *     case "windows":
 *         cmd = exec.Command("cmd", "/c", "start", "", path)
 *     }
 *     return cmd.Start()
 * }
 *
 * // ShowItemInFolder opens the folder containing the file and selects it.
 * func (s *Service) ShowItemInFolder(path string) error {
 *     var cmd *exec.Cmd
 *     switch runtime.GOOS {
 *     case "darwin":
 *         cmd = exec.Command("open", "-R", path)
 *     case "linux":
 *         cmd = exec.Command("xdg-open", filepath.Dir(path))
 *     case "windows":
 *         cmd = exec.Command("explorer", "/select,", path)
 *     }
 *     return cmd.Start()
 * }
 *
 * // TrashItem moves a file to the system trash.
 * func (s *Service) TrashItem(path string) error {
 *     switch runtime.GOOS {
 *     case "darwin":
 *         return exec.Command("osascript", "-e",
 *             `tell application "Finder" to delete POSIX file "`+path+`"`).Run()
 *     case "linux":
 *         return exec.Command("gio", "trash", path).Run()
 *     case "windows":
 *         // Windows requires PowerShell or COM for trash
 *         return exec.Command("powershell", "-Command",
 *             `Add-Type -AssemblyName Microsoft.VisualBasic; [Microsoft.VisualBasic.FileIO.FileSystem]::DeleteFile('`+path+`','OnlyErrorDialogs','SendToRecycleBin')`).Run()
 *     }
 *     return nil
 * }
 * ```
 */
