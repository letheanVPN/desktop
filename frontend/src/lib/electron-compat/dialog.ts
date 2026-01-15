/**
 * Electron dialog Compatibility Layer
 *
 * Maps Electron's dialog API to Wails Dialog system.
 *
 * Electron Concept -> Wails Equivalent:
 * - dialog.showOpenDialog()    -> Dialogs.Open()
 * - dialog.showSaveDialog()    -> Dialogs.Save()
 * - dialog.showMessageBox()    -> Dialogs.Info/Warning/Error/Question()
 *
 * @example
 * import { dialog } from '@lib/electron-compat';
 *
 * const result = await dialog.showOpenDialog({
 *   properties: ['openFile', 'multiSelections'],
 *   filters: [{ name: 'Images', extensions: ['jpg', 'png'] }]
 * });
 */

import { Dialogs } from '@wailsio/runtime';

export interface FileFilter {
  name: string;
  extensions: string[];
}

export interface OpenDialogOptions {
  title?: string;
  defaultPath?: string;
  buttonLabel?: string;
  filters?: FileFilter[];
  properties?: Array<
    | 'openFile'
    | 'openDirectory'
    | 'multiSelections'
    | 'showHiddenFiles'
    | 'createDirectory'
    | 'promptToCreate'
    | 'noResolveAliases'
    | 'treatPackageAsDirectory'
    | 'dontAddToRecent'
  >;
  message?: string;
}

export interface SaveDialogOptions {
  title?: string;
  defaultPath?: string;
  buttonLabel?: string;
  filters?: FileFilter[];
  message?: string;
  nameFieldLabel?: string;
  showsTagField?: boolean;
  properties?: Array<'showHiddenFiles' | 'createDirectory' | 'showOverwriteConfirmation' | 'dontAddToRecent'>;
}

export interface MessageBoxOptions {
  type?: 'none' | 'info' | 'error' | 'question' | 'warning';
  buttons?: string[];
  defaultId?: number;
  title?: string;
  message: string;
  detail?: string;
  checkboxLabel?: string;
  checkboxChecked?: boolean;
  cancelId?: number;
  noLink?: boolean;
}

export interface OpenDialogReturnValue {
  canceled: boolean;
  filePaths: string[];
}

export interface SaveDialogReturnValue {
  canceled: boolean;
  filePath?: string;
}

export interface MessageBoxReturnValue {
  response: number;
  checkboxChecked: boolean;
}

/**
 * Convert Electron file filters to Wails filter format
 */
function convertFilters(filters?: FileFilter[]): string {
  if (!filters || filters.length === 0) return '';

  // Wails uses pattern format: "*.jpg;*.png;*.gif"
  const patterns = filters.flatMap((f) => f.extensions.map((ext) => `*.${ext}`));
  return patterns.join(';');
}

export const dialog = {
  /**
   * Show a file open dialog.
   *
   * @param options - Dialog configuration options
   * @returns Promise resolving to selected file paths
   *
   * @example
   * const result = await dialog.showOpenDialog({
   *   title: 'Select Wallet File',
   *   filters: [{ name: 'Wallet', extensions: ['wallet', 'keys'] }],
   *   properties: ['openFile']
   * });
   *
   * if (!result.canceled) {
   *   console.log('Selected:', result.filePaths);
   * }
   */
  async showOpenDialog(options: OpenDialogOptions = {}): Promise<OpenDialogReturnValue> {
    const props = options.properties || ['openFile'];
    const isDirectory = props.includes('openDirectory');
    const allowMultiple = props.includes('multiSelections');

    try {
      let result: string | string[] | null;

      if (isDirectory) {
        // Wails directory selection
        result = await Dialogs.OpenDirectory({
          Title: options.title,
          DefaultDirectory: options.defaultPath,
          ButtonText: options.buttonLabel,
          CanCreateDirectories: props.includes('createDirectory'),
        });

        // Directory dialog returns single path or null
        return {
          canceled: !result,
          filePaths: result ? [result as string] : [],
        };
      } else {
        // Wails file selection
        if (allowMultiple) {
          result = await Dialogs.OpenMultipleFiles({
            Title: options.title,
            DefaultDirectory: options.defaultPath,
            DefaultFilename: '',
            ButtonText: options.buttonLabel,
            Filters: convertFilters(options.filters),
          });
        } else {
          result = await Dialogs.OpenFile({
            Title: options.title,
            DefaultDirectory: options.defaultPath,
            DefaultFilename: '',
            ButtonText: options.buttonLabel,
            Filters: convertFilters(options.filters),
          });
        }

        // Normalize to array
        const filePaths = result
          ? Array.isArray(result)
            ? result
            : [result]
          : [];

        return {
          canceled: filePaths.length === 0,
          filePaths,
        };
      }
    } catch (error) {
      console.error('[electron-compat] showOpenDialog error:', error);
      return { canceled: true, filePaths: [] };
    }
  },

  /**
   * Show a file save dialog.
   *
   * @param options - Dialog configuration options
   * @returns Promise resolving to the selected save path
   *
   * @example
   * const result = await dialog.showSaveDialog({
   *   title: 'Export Keys',
   *   defaultPath: 'my-wallet.keys',
   *   filters: [{ name: 'Keys', extensions: ['keys'] }]
   * });
   *
   * if (!result.canceled) {
   *   console.log('Saving to:', result.filePath);
   * }
   */
  async showSaveDialog(options: SaveDialogOptions = {}): Promise<SaveDialogReturnValue> {
    try {
      const result = await Dialogs.SaveFile({
        Title: options.title,
        DefaultDirectory: options.defaultPath ? options.defaultPath.split('/').slice(0, -1).join('/') : undefined,
        DefaultFilename: options.defaultPath ? options.defaultPath.split('/').pop() : undefined,
        ButtonText: options.buttonLabel,
        Filters: convertFilters(options.filters),
        CanCreateDirectories: options.properties?.includes('createDirectory'),
      });

      return {
        canceled: !result,
        filePath: result || undefined,
      };
    } catch (error) {
      console.error('[electron-compat] showSaveDialog error:', error);
      return { canceled: true };
    }
  },

  /**
   * Show a message box dialog.
   *
   * @param options - Message box configuration
   * @returns Promise resolving to the button index clicked
   *
   * @example
   * const result = await dialog.showMessageBox({
   *   type: 'question',
   *   buttons: ['Yes', 'No', 'Cancel'],
   *   title: 'Confirm',
   *   message: 'Are you sure you want to delete this wallet?',
   *   detail: 'This action cannot be undone.'
   * });
   *
   * if (result.response === 0) {
   *   // User clicked "Yes"
   * }
   */
  async showMessageBox(options: MessageBoxOptions): Promise<MessageBoxReturnValue> {
    try {
      // Map Electron dialog types to Wails dialog methods
      const dialogType = options.type || 'info';
      const buttons = options.buttons || ['OK'];

      // Wails has separate methods for each dialog type
      let dialogPromise: Promise<string>;

      const dialogOptions = {
        Title: options.title || '',
        Message: options.message,
        // Note: Wails dialogs don't support custom buttons in the same way
        // This is a simplified implementation
      };

      switch (dialogType) {
        case 'error':
          dialogPromise = Dialogs.Error(dialogOptions);
          break;
        case 'warning':
          dialogPromise = Dialogs.Warning(dialogOptions);
          break;
        case 'question':
          dialogPromise = Dialogs.Question(dialogOptions);
          break;
        case 'info':
        case 'none':
        default:
          dialogPromise = Dialogs.Info(dialogOptions);
          break;
      }

      const result = await dialogPromise;

      // Map Wails result to Electron-style response
      // Wails Question returns "Yes", "No", etc.
      const responseIndex = buttons.findIndex(
        (b) => b.toLowerCase() === (result || '').toLowerCase()
      );

      return {
        response: responseIndex >= 0 ? responseIndex : 0,
        checkboxChecked: false, // Wails doesn't support checkboxes in dialogs
      };
    } catch (error) {
      console.error('[electron-compat] showMessageBox error:', error);
      return { response: 0, checkboxChecked: false };
    }
  },

  /**
   * Show an error dialog (synchronous in Electron, async here).
   *
   * @param title - Dialog title
   * @param content - Error message content
   */
  async showErrorBox(title: string, content: string): Promise<void> {
    await Dialogs.Error({
      Title: title,
      Message: content,
    });
  },
};
