/**
 * Electron app Compatibility Layer
 *
 * Maps Electron's app API to Wails Application runtime.
 *
 * Note: Many Electron app APIs relate to the main process lifecycle,
 * which works differently in Wails. This provides the most commonly
 * used subset that makes sense in a Wails context.
 *
 * @example
 * import { app } from '@lib/electron-compat';
 *
 * console.log('App version:', app.getVersion());
 * console.log('User data path:', await app.getPath('userData'));
 */

import { Application, Call } from '@wailsio/runtime';

// Cache for app info to avoid repeated calls
let cachedAppInfo: { name: string; version: string } | null = null;

export const app = {
  /**
   * Get the application name.
   *
   * @returns The application name from wails.json
   */
  getName(): string {
    return 'Lethean Desktop'; // Could be made dynamic via Go binding
  },

  /**
   * Get the application version.
   *
   * @returns The application version
   *
   * @example
   * console.log(`Running version ${app.getVersion()}`);
   */
  getVersion(): string {
    // This could be bound from Go's build-time version
    return '1.0.0';
  },

  /**
   * Get a special directory path.
   *
   * NOTE: This requires a Go backend method to be implemented.
   *
   * @param name - The path type to get
   * @returns Promise resolving to the path string
   *
   * @example
   * const userDataPath = await app.getPath('userData');
   * const logsPath = await app.getPath('logs');
   */
  async getPath(
    name:
      | 'home'
      | 'appData'
      | 'userData'
      | 'sessionData'
      | 'temp'
      | 'exe'
      | 'module'
      | 'desktop'
      | 'documents'
      | 'downloads'
      | 'music'
      | 'pictures'
      | 'videos'
      | 'recent'
      | 'logs'
      | 'crashDumps'
  ): Promise<string> {
    try {
      // Maps to the config service's path resolution
      const result = await Call.ByName(
        'github.com/letheanVPN/desktop/services/core/config.Service.GetPath',
        name
      );
      return result as string;
    } catch {
      // Fallback to reasonable defaults
      console.warn(`[electron-compat] getPath('${name}') not implemented, using fallback`);
      return '';
    }
  },

  /**
   * Get the current application locale.
   *
   * @returns The system locale string (e.g., 'en-US')
   */
  getLocale(): string {
    return navigator.language || 'en-US';
  },

  /**
   * Get the system locale for spell checking.
   */
  getSystemLocale(): string {
    return navigator.language || 'en-US';
  },

  /**
   * Check if the app is packaged (production build).
   *
   * @returns true if running as packaged app
   */
  isPackaged(): boolean {
    // In Wails, check if we're in dev mode
    return !window.location.href.includes('localhost');
  },

  /**
   * Quit the application.
   *
   * @param exitCode - Optional exit code (default: 0)
   */
  quit(exitCode?: number): void {
    Application.Quit();
  },

  /**
   * Exit the application immediately.
   *
   * @param exitCode - Exit code (default: 0)
   */
  exit(exitCode?: number): void {
    Application.Quit();
  },

  /**
   * Relaunch the application.
   *
   * NOTE: Not directly supported in Wails - logs a warning.
   */
  relaunch(_options?: { args?: string[]; execPath?: string }): void {
    console.warn('[electron-compat] relaunch() is not directly supported in Wails');
    // Could potentially be implemented via Go with os/exec
  },

  /**
   * Check if the app is ready.
   * In Wails, the app is ready when the frontend loads.
   *
   * @returns Always true in the frontend context
   */
  isReady(): boolean {
    return true;
  },

  /**
   * Wait for the app to be ready.
   * Resolves immediately in Wails frontend context.
   *
   * @returns Promise that resolves when app is ready
   */
  whenReady(): Promise<void> {
    return Promise.resolve();
  },

  /**
   * Focus the application.
   */
  focus(_options?: { steal: boolean }): void {
    window.focus();
  },

  /**
   * Hide the application (macOS).
   * Maps to Window.Hide() in Wails.
   */
  hide(): void {
    // Would need Go binding for proper implementation
    console.warn('[electron-compat] hide() requires Go backend implementation');
  },

  /**
   * Show the application (after hide).
   */
  show(): void {
    window.focus();
  },

  /**
   * Set the application badge count (macOS/Linux).
   *
   * @param count - Badge count (0 to clear)
   * @returns Whether the call succeeded
   */
  setBadgeCount(count: number): boolean {
    // Not directly supported in Wails - would need Go implementation
    console.warn('[electron-compat] setBadgeCount() requires Go backend implementation');
    return false;
  },

  /**
   * Get the badge count.
   *
   * @returns The current badge count
   */
  getBadgeCount(): number {
    return 0;
  },

  /**
   * Check if running on Rosetta 2 (Apple Silicon with x64 binary).
   *
   * @returns Promise resolving to boolean
   */
  async isRunningUnderARM64Translation(): Promise<boolean> {
    // Would need Go implementation to check
    return false;
  },

  // =========================================================================
  // Event-like methods (for Electron compatibility)
  // =========================================================================

  /**
   * Register a callback for when the app is ready.
   * Executes immediately since Wails frontend is already "ready".
   *
   * @param callback - Function to call
   */
  on(event: string, callback: (...args: unknown[]) => void): void {
    if (event === 'ready') {
      // Already ready in frontend context
      callback();
    } else if (event === 'window-all-closed') {
      // Not applicable - Wails handles this
      console.warn(`[electron-compat] app.on('${event}') - event not supported in Wails frontend`);
    } else if (event === 'activate') {
      // macOS dock click - would need Go implementation
      console.warn(`[electron-compat] app.on('${event}') requires Go backend implementation`);
    } else {
      console.warn(`[electron-compat] app.on('${event}') - unknown event`);
    }
  },

  /**
   * Register a one-time callback.
   */
  once(event: string, callback: (...args: unknown[]) => void): void {
    this.on(event, callback);
  },
};

/*
 * =============================================================================
 * GO BACKEND IMPLEMENTATION REQUIRED (for getPath)
 * =============================================================================
 *
 * Add this method to: services/core/config/service.go
 *
 * ```go
 * import (
 *     "os"
 *     "path/filepath"
 *     "runtime"
 *
 *     "github.com/adrg/xdg"
 * )
 *
 * // GetPath returns special directory paths (Electron app.getPath compatibility)
 * func (s *Service) GetPath(name string) (string, error) {
 *     switch name {
 *     case "home":
 *         return os.UserHomeDir()
 *     case "appData":
 *         return xdg.ConfigHome, nil
 *     case "userData":
 *         return xdg.DataHome + "/lethean-desktop", nil
 *     case "temp":
 *         return os.TempDir(), nil
 *     case "desktop":
 *         home, _ := os.UserHomeDir()
 *         return filepath.Join(home, "Desktop"), nil
 *     case "documents":
 *         home, _ := os.UserHomeDir()
 *         return filepath.Join(home, "Documents"), nil
 *     case "downloads":
 *         home, _ := os.UserHomeDir()
 *         return filepath.Join(home, "Downloads"), nil
 *     case "logs":
 *         return xdg.StateHome + "/lethean-desktop/logs", nil
 *     default:
 *         return "", fmt.Errorf("unknown path name: %s", name)
 *     }
 * }
 * ```
 */
