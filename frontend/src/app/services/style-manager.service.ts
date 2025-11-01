import { Injectable } from '@angular/core';
import { environment } from '../../environments/environment';

@Injectable({
  providedIn: 'root'
})
export class StyleManagerService {

  constructor() { }

  init() {
    // if (environment.pro) {
    //   this.loadProStyles();
    // }
  }

  private loadProStyles() {
    const proStyles = [
      'node_modules/@fortawesome/fontawesome-free/css/light.min.css',
      'node_modules/@fortawesome/fontawesome-free/css/jelly-regular.min.css',
      'node_modules/@fortawesome/fontawesome-free/css/jelly-fill-regular.min.css'
    ];

    proStyles.forEach(style => {
      const link = document.createElement('link');
      link.rel = 'stylesheet';
      link.href = style;
      document.head.appendChild(link);
    });
  }
}
