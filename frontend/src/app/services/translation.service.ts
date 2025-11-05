import { Injectable, isDevMode } from '@angular/core';
import { GetAllMessages, Translate } from '@lthn/core/i18n/service';
import { TranslateService } from '@ngx-translate/core';
import { firstValueFrom } from 'rxjs';

@Injectable({
  providedIn: 'root'
})
export class TranslationService {
  private translations: Map<string, string> = new Map();
  private isLoaded = false;
  private loadingPromise: Promise<void>;

  constructor(private ngxTranslate: TranslateService) {
    this.loadingPromise = this.loadTranslations('en');
  }

  public reload(lang: string): Promise<void> {
    this.isLoaded = false;
    this.loadingPromise = this.loadTranslations(lang);
    return this.loadingPromise;
  }

  private async loadTranslations(lang: string): Promise<void> {
    if (isDevMode()) {
      await firstValueFrom(this.ngxTranslate.use(lang));
      this.isLoaded = true;
      console.log('TranslationService: Using ngx-translate for development.');
    } else {
      try {
        const allMessages: Record<string, string> | null = await GetAllMessages(lang);
        if (!allMessages) {
          return
        }
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
  }

  public translate(key: string): string {
    if (isDevMode()) {
      return this.ngxTranslate.instant(key);
    } else {
      if (!this.isLoaded) {
        return key;
      }
      return this.translations.get(key) || key;
    }
  }

  public _ = this.translate;

  public async translateOnDemand(key: string): Promise<string> {
    if (isDevMode()) {
      return firstValueFrom(this.ngxTranslate.get(key));
    } else {
      try {
        return await Translate(key).then(s => s || key);
      } catch (error) {
        console.error(`TranslationService: Failed to translate key "${key}" on demand:`, error);
        return key; // Fallback
      }
    }
  }

  public onReady(): Promise<void> {
    return this.loadingPromise;
  }
}
