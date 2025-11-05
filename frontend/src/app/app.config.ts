import {APP_INITIALIZER, ApplicationConfig, importProvidersFrom, isDevMode} from '@angular/core';
import {provideRouter, withHashLocation} from '@angular/router';
import {MonacoEditorModule} from 'ngx-monaco-editor-v2';

import {routes} from './app.routes';
import {I18nService} from './services/i18n.service';
import {TranslateModule} from '@ngx-translate/core';
import {provideHttpClient} from '@angular/common/http';
import {provideTranslateHttpLoader} from '@ngx-translate/http-loader';
import {StyleManagerService} from './services/style-manager.service';
import { provideHighcharts } from 'highcharts-angular';

const translationProviders = [
  provideHttpClient(),
  importProvidersFrom(
    TranslateModule.forRoot({
      fallbackLang: 'en',
    })
  ),
  provideHighcharts({
    options: {
      credits: {enabled: false},
      title: {
        style: {
          color: 'tomato',
        },
      },
      legend: {
        enabled: false,
      },
    },
    modules: () => {
      return [
        import('highcharts/esm/modules/accessibility'),
        import('highcharts/esm/modules/exporting'),
        import('highcharts/esm/themes/sunset'),
      ];
    },
  }),
  ...(isDevMode()
    ? [
      provideTranslateHttpLoader({
        prefix: './assets/i18n/',
        suffix: '.json',
      }),
    ]
    : []),
];

export function initializeApp(styleManager: StyleManagerService) {
  return () => styleManager.init();
}

export const appConfig: ApplicationConfig = {
  providers: [
    provideRouter(routes, withHashLocation()),
    importProvidersFrom(MonacoEditorModule.forRoot()),
    I18nService,
    ...translationProviders,
    {
      provide: APP_INITIALIZER,
      useFactory: initializeApp,
      deps: [StyleManagerService],
      multi: true
    }
  ],
};
