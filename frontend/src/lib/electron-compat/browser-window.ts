/**
 * Electron BrowserWindow Compatibility Layer
 *
 * Maps Electron's BrowserWindow API to Wails Window system.
 *
 * IMPORTANT: Wails has a fundamentally different window model than Electron:
 * - Electron: Multiple BrowserWindow instances, each with own renderer process
 * - Wails: Single main window with ability to spawn additional windows via Go
 *
 * This compatibility layer provides a subset of BrowserWindow functionality
 * that maps to Wails' window capabilities.
 *
 * @example
 * import { BrowserWindow } from '@lib/electron-compat';
 *
 * // Get current window
 * const win = BrowserWindow.getFocusedWindow();
 * win.setTitle('My Window');
 * win.maximize();
 */

import { Window, Call, Events } from '@wailsio/runtime';

export interface BrowserWindowOptions {
  width?: number;
  height?: number;
  x?: number;
  y?: number;
  minWidth?: number;
  minHeight?: number;
  maxWidth?: number;
  maxHeight?: number;
  resizable?: boolean;
  movable?: boolean;
  minimizable?: boolean;
  maximizable?: boolean;
  closable?: boolean;
  focusable?: boolean;
  alwaysOnTop?: boolean;
  fullscreen?: boolean;
  fullscreenable?: boolean;
  title?: string;
  show?: boolean;
  frame?: boolean;
  transparent?: boolean;
  backgroundColor?: string;
}

export interface Rectangle {
  x: number;
  y: number;
  width: number;
  height: number;
}

/**
 * BrowserWindow compatibility class for the current Wails window.
 *
 * Note: Unlike Electron, you cannot create new BrowserWindow instances
 * directly from the frontend. Use the Go backend's display service
 * to open new windows.
 */
export class BrowserWindow {
  private id: number;
  private static currentWindow: BrowserWindow | null = null;

  private constructor(id: number = 0) {
    this.id = id;
  }

  /**
   * Get the currently focused window.
   * In Wails, this typically returns a wrapper for the main window.
   *
   * @returns BrowserWindow instance or null
   */
  static getFocusedWindow(): BrowserWindow | null {
    if (!BrowserWindow.currentWindow) {
      BrowserWindow.currentWindow = new BrowserWindow(0);
    }
    return BrowserWindow.currentWindow;
  }

  /**
   * Get all open windows.
   * In Wails, window management is handled by Go, so this returns
   * just the current window context.
   *
   * @returns Array of BrowserWindow instances
   */
  static getAllWindows(): BrowserWindow[] {
    const focused = BrowserWindow.getFocusedWindow();
    return focused ? [focused] : [];
  }

  /**
   * Create a new browser window.
   *
   * NOTE: In Wails, new windows must be created via Go backend.
   * This method calls the display service to open a new window.
   *
   * @param options - Window configuration
   */
  static async create(options: BrowserWindowOptions & { url?: string; name?: string }): Promise<void> {
    try {
      await Call.ByName(
        'github.com/letheanVPN/desktop/services/core/display.Service.OpenWindow',
        options.name || 'window',
        {
          Title: options.title || '',
          Width: options.width || 800,
          Height: options.height || 600,
          URL: options.url || '/',
          AlwaysOnTop: options.alwaysOnTop || false,
          Frameless: options.frame === false,
          Resizable: options.resizable !== false,
          MinWidth: options.minWidth,
          MinHeight: options.minHeight,
          MaxWidth: options.maxWidth,
          MaxHeight: options.maxHeight,
        }
      );
    } catch (error) {
      console.error('[electron-compat] BrowserWindow.create failed:', error);
      throw error;
    }
  }

  // =========================================================================
  // Instance Methods - Window State
  // =========================================================================

  /**
   * Close the window.
   */
  close(): void {
    Window.Close();
  }

  /**
   * Focus the window.
   */
  focus(): void {
    Window.Focus();
  }

  /**
   * Blur (unfocus) the window.
   */
  blur(): void {
    // Not directly supported in Wails
    console.warn('[electron-compat] blur() not directly supported');
  }

  /**
   * Check if the window is focused.
   */
  isFocused(): boolean {
    return document.hasFocus();
  }

  /**
   * Check if the window is destroyed/closed.
   */
  isDestroyed(): boolean {
    return false; // Current window is never destroyed in this context
  }

  /**
   * Show the window.
   */
  show(): void {
    Window.Show();
  }

  /**
   * Hide the window.
   */
  hide(): void {
    Window.Hide();
  }

  /**
   * Check if the window is visible.
   */
  isVisible(): boolean {
    return !document.hidden;
  }

  /**
   * Check if the window is maximized.
   */
  async isMaximized(): Promise<boolean> {
    return await Window.IsMaximised();
  }

  /**
   * Maximize the window.
   */
  maximize(): void {
    Window.Maximise();
  }

  /**
   * Unmaximize the window.
   */
  unmaximize(): void {
    Window.UnMaximise();
  }

  /**
   * Check if the window is minimized.
   */
  async isMinimized(): Promise<boolean> {
    return await Window.IsMinimised();
  }

  /**
   * Minimize the window.
   */
  minimize(): void {
    Window.Minimise();
  }

  /**
   * Restore the window from minimized state.
   */
  restore(): void {
    Window.UnMinimise();
  }

  /**
   * Check if the window is in fullscreen mode.
   */
  async isFullScreen(): Promise<boolean> {
    return await Window.IsFullscreen();
  }

  /**
   * Set fullscreen mode.
   */
  setFullScreen(flag: boolean): void {
    if (flag) {
      Window.Fullscreen();
    } else {
      Window.UnFullscreen();
    }
  }

  /**
   * Toggle fullscreen mode.
   */
  toggleFullScreen(): void {
    Window.ToggleFullscreen();
  }

  // =========================================================================
  // Instance Methods - Window Properties
  // =========================================================================

  /**
   * Get the window title.
   */
  getTitle(): string {
    return document.title;
  }

  /**
   * Set the window title.
   */
  setTitle(title: string): void {
    Window.SetTitle(title);
  }

  /**
   * Get the window bounds.
   */
  async getBounds(): Promise<Rectangle> {
    const size = await Window.Size();
    const pos = await Window.Position();
    return {
      x: pos.x,
      y: pos.y,
      width: size.width,
      height: size.height,
    };
  }

  /**
   * Set the window bounds.
   */
  setBounds(bounds: Partial<Rectangle>): void {
    if (bounds.width !== undefined && bounds.height !== undefined) {
      Window.SetSize(bounds.width, bounds.height);
    }
    if (bounds.x !== undefined && bounds.y !== undefined) {
      Window.SetPosition(bounds.x, bounds.y);
    }
  }

  /**
   * Get the window size.
   */
  async getSize(): Promise<[number, number]> {
    const size = await Window.Size();
    return [size.width, size.height];
  }

  /**
   * Set the window size.
   */
  setSize(width: number, height: number): void {
    Window.SetSize(width, height);
  }

  /**
   * Get the window position.
   */
  async getPosition(): Promise<[number, number]> {
    const pos = await Window.Position();
    return [pos.x, pos.y];
  }

  /**
   * Set the window position.
   */
  setPosition(x: number, y: number): void {
    Window.SetPosition(x, y);
  }

  /**
   * Center the window on screen.
   */
  center(): void {
    Window.Center();
  }

  /**
   * Set minimum window size.
   */
  setMinimumSize(width: number, height: number): void {
    Window.SetMinSize(width, height);
  }

  /**
   * Set maximum window size.
   */
  setMaximumSize(width: number, height: number): void {
    Window.SetMaxSize(width, height);
  }

  /**
   * Set whether the window is resizable.
   */
  setResizable(resizable: boolean): void {
    Window.SetResizable(resizable);
  }

  /**
   * Check if the window is resizable.
   */
  async isResizable(): Promise<boolean> {
    return await Window.Resizable();
  }

  /**
   * Set always on top.
   */
  setAlwaysOnTop(flag: boolean): void {
    Window.SetAlwaysOnTop(flag);
  }

  /**
   * Check if always on top.
   */
  async isAlwaysOnTop(): Promise<boolean> {
    return await Window.IsAlwaysOnTop();
  }

  /**
   * Set the window background color.
   */
  setBackgroundColor(_color: string): void {
    // Simplified - would need color parsing
    Window.SetBackgroundColour({ r: 0, g: 0, b: 0, a: 255 });
  }

  // =========================================================================
  // Instance Methods - Events
  // =========================================================================

  /**
   * Register an event listener.
   */
  on(event: string, listener: (...args: unknown[]) => void): this {
    const eventMap: Record<string, string> = {
      close: 'window:close',
      closed: 'window:closed',
      focus: 'window:focus',
      blur: 'window:blur',
      maximize: 'window:maximise',
      unmaximize: 'window:unmaximise',
      minimize: 'window:minimise',
      restore: 'window:restore',
      resize: 'window:resize',
      move: 'window:move',
      'enter-full-screen': 'window:fullscreen',
      'leave-full-screen': 'window:unfullscreen',
    };

    const wailsEvent = eventMap[event];
    if (wailsEvent) {
      Events.On(wailsEvent, listener);
    } else {
      console.warn(`[electron-compat] BrowserWindow event '${event}' not mapped`);
    }

    return this;
  }

  /**
   * Register a one-time event listener.
   */
  once(event: string, listener: (...args: unknown[]) => void): this {
    const eventMap: Record<string, string> = {
      close: 'window:close',
      ready: 'window:ready',
    };

    const wailsEvent = eventMap[event];
    if (wailsEvent) {
      Events.Once(wailsEvent, listener);
    }

    return this;
  }

  // =========================================================================
  // WebContents-like methods (limited support)
  // =========================================================================

  /**
   * Get the webContents-like object.
   * Returns a simplified interface since Wails doesn't have webContents.
   */
  get webContents() {
    return {
      /**
       * Get the current URL.
       */
      getURL: (): string => {
        return window.location.href;
      },

      /**
       * Navigate to a URL (changes the hash route).
       */
      loadURL: (url: string): void => {
        if (url.startsWith('#')) {
          window.location.hash = url;
        } else {
          window.location.href = url;
        }
      },

      /**
       * Reload the page.
       */
      reload: (): void => {
        window.location.reload();
      },

      /**
       * Open DevTools.
       */
      openDevTools: (): void => {
        console.log('[electron-compat] To open DevTools, use browser developer tools (F12 or Cmd+Option+I)');
      },

      /**
       * Send a message to the renderer (no-op in Wails, we ARE the renderer).
       */
      send: (channel: string, ...args: unknown[]): void => {
        Events.Emit({ name: channel, data: args });
      },
    };
  }
}

// Also export as default for compatibility with some Electron patterns
export default BrowserWindow;
