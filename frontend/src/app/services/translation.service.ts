import { Injectable } from '@angular/core';
import { GetAllMessages, Translate } from '@lthn/core/i18n/service';

@Injectable({
  providedIn: 'root'
})
export class TranslationService {
  private translations: Map<string, string> = new Map();
  private isLoaded = false;
  private loadingPromise: Promise<void>;

  constructor() {
    // Load initial translations
    this.loadingPromise = this.loadTranslations();
  }

  /**
   * Reloads translations from the backend.
   * This is called by I18nService when the language changes.
   */
  public reload(): Promise<void> {
      this.isLoaded = false;
      this.loadingPromise = this.loadTranslations();
      return this.loadingPromise;
  }

  private async loadTranslations(): Promise<void> {
    try {
      const allMessages: Record<string, string> = await GetAllMessages('en');
      this.translations.clear();
      for (const key in allMessages) {
        if (Object.prototype.hasOwnProperty.call(allMessages, key)) {
          this.translations.set(key, allMessages[key]);
        }
      }
      this.isLoaded = true;
      console.log('TranslationService: Translations loaded/reloaded successfully.');
    } catch (error) {
      console.error('TranslationService: Failed to load translations:', error);
      throw error;
    }
  }

  /**
   * Synchronously translates a message key.
   * @param key The translation key.
   * @returns The translated string, or the key if not found.
   */
  public translate(key: string): string {
    if (!this.isLoaded) {
      return key;
    }
    return this.translations.get(key) || key;
  }

  /**
   * Alias for the synchronous translate method, for more concise template usage.
   */
  public _ = this.translate;

  /**
   * Asynchronously translates a single key on-demand by calling the backend.
   * @param key The translation key.
   * @returns A promise that resolves with the translated string, or the key on error.
   */
  public async translateOnDemand(key: string): Promise<string> {
    try {
      return await Translate(key).then(s => s || key);
    } catch (error) {
      console.error(`TranslationService: Failed to translate key "${key}" on demand:`, error);
      return key; // Fallback
    }
  }

  /**
   * Returns a promise that resolves when the current translation load is complete.
   */
  public onReady(): Promise<void> {
    return this.loadingPromise;
  }
}
