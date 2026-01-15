/**
 * Electron Compatibility Layer for Wails v3
 *
 * This module provides Electron-like APIs that map to Wails runtime equivalents.
 * It's designed to help developers familiar with Electron contribute to this
 * Wails-based application without needing to learn the Wails API from scratch.
 *
 * Usage:
 *   import { ipcRenderer, shell, dialog, app } from '@lib/electron-compat';
 *
 *   // Works like Electron!
 *   ipcRenderer.invoke('my-channel', data);
 *   shell.openExternal('https://lethean.io');
 *
 * @see https://wails.io/docs/reference/runtime/intro
 * @see https://www.electronjs.org/docs/latest/api/ipc-renderer
 */

export { ipcRenderer } from './ipc-renderer';
export { shell } from './shell';
export { dialog } from './dialog';
export { app } from './app';
export { clipboard } from './clipboard';
export { BrowserWindow } from './browser-window';

// Re-export types for TypeScript users
export type { IpcRendererEvent, IpcMainInvokeEvent } from './ipc-renderer';
export type { OpenDialogOptions, SaveDialogOptions, MessageBoxOptions } from './dialog';
