/**
 * Electron clipboard Compatibility Layer
 *
 * Maps Electron's clipboard API to the browser Clipboard API.
 *
 * The browser's Clipboard API is async and requires user permissions,
 * while Electron's is sync. This implementation uses async methods
 * with fallbacks where possible.
 *
 * @example
 * import { clipboard } from '@lib/electron-compat';
 *
 * await clipboard.writeText('Hello, World!');
 * const text = await clipboard.readText();
 */

export const clipboard = {
  /**
   * Read plain text from the clipboard.
   *
   * @param type - Clipboard type ('selection' for Linux primary selection)
   * @returns Promise resolving to clipboard text content
   *
   * @example
   * const walletAddress = await clipboard.readText();
   */
  async readText(type?: 'selection' | 'clipboard'): Promise<string> {
    try {
      return await navigator.clipboard.readText();
    } catch (error) {
      console.error('[electron-compat] clipboard.readText failed:', error);
      return '';
    }
  },

  /**
   * Write plain text to the clipboard.
   *
   * @param text - The text to write
   * @param type - Clipboard type ('selection' for Linux primary selection)
   *
   * @example
   * await clipboard.writeText(walletAddress);
   */
  async writeText(text: string, type?: 'selection' | 'clipboard'): Promise<void> {
    try {
      await navigator.clipboard.writeText(text);
    } catch (error) {
      // Fallback for older browsers or when clipboard API is blocked
      console.error('[electron-compat] clipboard.writeText failed:', error);
      fallbackCopyText(text);
    }
  },

  /**
   * Read HTML content from the clipboard.
   *
   * @returns Promise resolving to HTML string
   */
  async readHTML(): Promise<string> {
    try {
      const items = await navigator.clipboard.read();
      for (const item of items) {
        if (item.types.includes('text/html')) {
          const blob = await item.getType('text/html');
          return await blob.text();
        }
      }
      return '';
    } catch (error) {
      console.error('[electron-compat] clipboard.readHTML failed:', error);
      return '';
    }
  },

  /**
   * Write HTML content to the clipboard.
   *
   * @param markup - The HTML string to write
   */
  async writeHTML(markup: string): Promise<void> {
    try {
      const blob = new Blob([markup], { type: 'text/html' });
      await navigator.clipboard.write([
        new ClipboardItem({
          'text/html': blob,
        }),
      ]);
    } catch (error) {
      console.error('[electron-compat] clipboard.writeHTML failed:', error);
    }
  },

  /**
   * Read an image from the clipboard.
   *
   * @returns Promise resolving to image data as data URL, or empty string
   */
  async readImage(): Promise<string> {
    try {
      const items = await navigator.clipboard.read();
      for (const item of items) {
        const imageTypes = item.types.filter((type) => type.startsWith('image/'));
        if (imageTypes.length > 0) {
          const blob = await item.getType(imageTypes[0]);
          return await blobToDataURL(blob);
        }
      }
      return '';
    } catch (error) {
      console.error('[electron-compat] clipboard.readImage failed:', error);
      return '';
    }
  },

  /**
   * Write an image to the clipboard.
   *
   * @param dataUrl - Image as data URL (e.g., 'data:image/png;base64,...')
   */
  async writeImage(dataUrl: string): Promise<void> {
    try {
      const response = await fetch(dataUrl);
      const blob = await response.blob();
      await navigator.clipboard.write([
        new ClipboardItem({
          [blob.type]: blob,
        }),
      ]);
    } catch (error) {
      console.error('[electron-compat] clipboard.writeImage failed:', error);
    }
  },

  /**
   * Check if the clipboard has content of a specific format.
   *
   * @param format - MIME type to check for
   * @returns Promise resolving to boolean
   */
  async has(format: string): Promise<boolean> {
    try {
      const items = await navigator.clipboard.read();
      return items.some((item) => item.types.includes(format));
    } catch {
      return false;
    }
  },

  /**
   * Clear the clipboard.
   */
  async clear(): Promise<void> {
    try {
      await navigator.clipboard.writeText('');
    } catch (error) {
      console.error('[electron-compat] clipboard.clear failed:', error);
    }
  },

  /**
   * Get available formats in the clipboard.
   *
   * @returns Promise resolving to array of MIME types
   */
  async availableFormats(): Promise<string[]> {
    try {
      const items = await navigator.clipboard.read();
      const formats: string[] = [];
      for (const item of items) {
        formats.push(...item.types);
      }
      return [...new Set(formats)];
    } catch {
      return [];
    }
  },

  // =========================================================================
  // Electron sync methods (not supported - use async versions above)
  // =========================================================================

  /**
   * @deprecated Use readText() instead - sync clipboard not supported in browser
   */
  readTextSync(): string {
    console.warn('[electron-compat] readTextSync not supported. Use async readText() instead.');
    return '';
  },

  /**
   * @deprecated Use writeText() instead - sync clipboard not supported in browser
   */
  writeTextSync(_text: string): void {
    console.warn('[electron-compat] writeTextSync not supported. Use async writeText() instead.');
  },
};

/**
 * Fallback copy using execCommand (for older browsers).
 */
function fallbackCopyText(text: string): void {
  const textArea = document.createElement('textarea');
  textArea.value = text;
  textArea.style.position = 'fixed';
  textArea.style.left = '-9999px';
  textArea.style.top = '-9999px';
  document.body.appendChild(textArea);
  textArea.focus();
  textArea.select();

  try {
    document.execCommand('copy');
  } catch (err) {
    console.error('[electron-compat] Fallback copy failed:', err);
  }

  document.body.removeChild(textArea);
}

/**
 * Convert a Blob to a data URL.
 */
function blobToDataURL(blob: Blob): Promise<string> {
  return new Promise((resolve, reject) => {
    const reader = new FileReader();
    reader.onloadend = () => resolve(reader.result as string);
    reader.onerror = reject;
    reader.readAsDataURL(blob);
  });
}
