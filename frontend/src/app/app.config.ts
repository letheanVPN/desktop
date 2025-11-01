import { ApplicationConfig, importProvidersFrom, isDevMode } from '@angular/core';
import { provideRouter, withHashLocation } from '@angular/router';
import { MonacoEditorModule } from 'ngx-monaco-editor-v2';

import { routes } from './app.routes';
import { I18nService } from './services/i18n.service';
import { TranslateModule } from '@ngx-translate/core';
import { provideHttpClient } from '@angular/common/http';
import { provideTranslateHttpLoader } from '@ngx-translate/http-loader';

const translationProviders = isDevMode()
  ? [
      provideHttpClient(),
      importProvidersFrom(
        TranslateModule.forRoot({
          fallbackLang: 'en',
        })
      ),
      provideTranslateHttpLoader({
        prefix: './assets/i18n/',
        suffix: '.json',
      }),
    ]
  : [];

export const appConfig: ApplicationConfig = {
  providers: [
    provideRouter(routes, withHashLocation()),
    importProvidersFrom(MonacoEditorModule.forRoot()),
    I18nService,
    ...translationProviders,
  ],
};
