import { Injectable, isDevMode, Optional } from '@angular/core';
import { SetLanguage, AvailableLanguages } from '@lthn/core/i18n/service';
import { BehaviorSubject } from 'rxjs';
import { TranslationService } from './translation.service';
import { TranslateService } from '@ngx-translate/core';

@Injectable({
  providedIn: 'root'
})
export class I18nService {
  private currentLanguageSubject = new BehaviorSubject<string>('en');
  public currentLanguage$ = this.currentLanguageSubject.asObservable();

  constructor(
    private translationService: TranslationService,
    @Optional() private ngxTranslate?: TranslateService
  ) {
    if (isDevMode() && this.ngxTranslate) {
      this.ngxTranslate.setDefaultLang('en');
    }
  }

  async setLanguage(lang: string): Promise<void> {
    if (isDevMode() && this.ngxTranslate) {
      await this.translationService.reload(lang);
      this.currentLanguageSubject.next(lang);
    } else {
      try {
        await SetLanguage(lang);
        this.currentLanguageSubject.next(lang);
        await this.translationService.reload(lang);
      } catch (error) {
        console.error(`I18nService: Failed to set language to "${lang}":`, error);
        throw error;
      }
    }
  }

  getAvailableLanguages(): Promise<string[]> {
    if (isDevMode()) {
      return Promise.resolve(['en']); // For dev, we can mock this.
    }
    return AvailableLanguages().then((languages) => {
      if (languages == null || languages?.length == 0)
        return Promise.resolve(['en']);
      return Promise.resolve(languages)
    });
  }

  public onReady(): Promise<void> {
    return this.translationService.onReady();
  }
}
