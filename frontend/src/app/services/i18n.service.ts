import { Injectable } from '@angular/core';
import { SetLanguage, AvailableLanguages } from '@lthn/core/i18n/service';
import { BehaviorSubject } from 'rxjs';
import { TranslationService } from './translation.service';

@Injectable({
  providedIn: 'root'
})
export class I18nService {
  private currentLanguageSubject = new BehaviorSubject<string>('en');
  public currentLanguage$ = this.currentLanguageSubject.asObservable();

  constructor(private translationService: TranslationService) {
    // TranslationService handles the initial load.
    // This service will configure other aspects like date/number formatting.
  }

  /**
   * Asynchronously sets the application language in the backend and reloads all translations.
   * @param lang The language code (e.g., "en", "es").
   */
  async setLanguage(lang: string): Promise<void> {
    try {
      await SetLanguage(lang);
      this.currentLanguageSubject.next(lang);
      // Tell TranslationService to reload translations for the new language.
      await this.translationService.reload();
    } catch (error) {
      console.error(`I18nService: Failed to set language to "${lang}":`, error);
      throw error;
    }
  }

  /**
   * Returns the list of available languages from the backend.
   * @returns A Promise that resolves with an array of language tags.
   */
  getAvailableLanguages(): Promise<string[]> {
    return AvailableLanguages();
  }

  /**
   * Returns a promise that resolves when the translation load is complete.
   * This can be used in components to wait for translations before rendering.
   */
  public onReady(): Promise<void> {
    return this.translationService.onReady();
  }
}
