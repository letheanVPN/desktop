import { Component } from '@angular/core';
import { CommonModule } from '@angular/common';

@Component({
  selector: 'system-tray-frame',
  standalone: true,
  imports: [CommonModule],
  template: `
    <div class="flex flex-col h-screen overflow-hidden rounded-md bg-white shadow-sm dark:bg-gray-800/50 dark:shadow-none dark:outline dark:-outline-offset-0 dark:outline-white/10">
      <div class="flex items-center justify-between px-4 py-4 bg-gray-50 dark:bg-gray-800 border-b border-gray-200 dark:border-white/10">
        <div class="flex h-8 shrink-0 items-center">
          <img class="h-58 w-80 pt-1" src="./logo/lthn/logo-full-gradient.png" alt="Lethean Community">
        </div>
        <div class="relative">
          <button (click)="settingsMenuOpen = !settingsMenuOpen" class="relative flex items-center">
            <span class="absolute -inset-1.5"></span>
            <span class="sr-only">Open settings menu</span>
            <svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" stroke-width="1.5" stroke="currentColor" class="w-6 h-6 text-gray-400">
              <path stroke-linecap="round" stroke-linejoin="round" d="M9.594 3.94c.09-.542.56-.94 1.11-.94h2.593c.55 0 1.02.398 1.11.94l.213 1.281c.063.374.313.686.645.87.074.04.147.083.22.127.324.196.72.257 1.075.124l1.217-.456a1.125 1.125 0 011.37.49l1.296 2.247a1.125 1.125 0 01-.26 1.431l-1.003.827c-.293.24-.438.613-.438.995s.145.755.438.995l1.003.827c.424.35.534.954.26 1.431l-1.296 2.247a1.125 1.125 0 01-1.37.49l-1.217-.456c-.355-.133-.75-.072-1.076.124a6.57 6.57 0 01-.22.128c-.331.183-.581.495-.644.869l-.213 1.28c-.09.543-.56.941-1.11.941h-2.594c-.55 0-1.02-.398-1.11-.94l-.213-1.281c-.063-.374-.313-.686-.645-.87a6.52 6.52 0 01-.22-.127c-.324-.196-.72-.257-1.075-.124l-1.217.456a1.125 1.125 0 01-1.37-.49l-1.296-2.247a1.125 1.125 0 01.26-1.431l1.004-.827c.292-.24.437-.613.437-.995s-.145-.755-.437-.995l-1.004-.827a1.125 1.125 0 01-.26-1.431l1.296-2.247a1.125 1.125 0 011.37-.49l1.217.456c.355.133.75.072 1.076-.124.072-.044.146-.087.22-.128.332-.183.582-.495.644-.869l.213-1.28z" />
              <path stroke-linecap="round" stroke-linejoin="round" d="M15 12a3 3 0 11-6 0 3 3 0 016 0z" />
            </svg>
          </button>
          <div *ngIf="settingsMenuOpen" class="absolute right-0 z-10 mt-2.5 w-48 origin-top-right rounded-md bg-white py-2 shadow-lg outline outline-gray-900/5 transition-discrete [--anchor-gap:--spacing(2.5)] data-closed:scale-95 data-closed:transform data-closed:opacity-0 data-enter:duration-100 data-enter:ease-out data-leave:duration-75 data-leave:ease-in dark:bg-gray-800 dark:shadow-none dark:-outline-offset-1 dark:outline-white/10">
            <a *ngFor="let item of settingsNavigation" [href]="item.href" class="block px-3 py-1 text-sm/6 text-gray-900 focus:bg-gray-50 focus:outline-hidden dark:text-white dark:focus:bg-white/5">{{ item.name }}</a>
          </div>
        </div>
      </div>
      <div class="flex-grow bg-gray-50 overflow-y-auto">
        <ul role="list" class="divide-y divide-gray-200 dark:divide-white/10">
          <li class="px-6 py-4">
            <p>Status: Connected</p>
          </li>
          <li class="px-6 py-4">
            <p>IP: 127.0.0.1</p>
          </li>
          <li class="px-6 py-4">
            <p>Uptime: 00:00:00</p>
          </li>
        </ul>
      </div>
      <div class="px-6 py-4 bg-gray-50 dark:bg-gray-800 border-t border-gray-200 dark:border-white/10">
        <button class="text-sm text-red-600 dark:text-red-400">Quit</button>
      </div>
    </div>
  `
})
export class SystemTrayFrame {
  settingsMenuOpen = false;

  settingsNavigation = [
    { name: 'Settings', href: '#' },
    { name: 'About', href: '#' },
    { name: 'Check for Updates...', href: '#' },
  ];
}
