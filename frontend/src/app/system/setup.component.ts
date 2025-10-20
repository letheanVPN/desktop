import { Component, OnInit, OnDestroy } from '@angular/core';
import { CommonModule } from '@angular/common';
import { RouterLink, RouterOutlet, Router, ActivatedRoute, NavigationEnd } from '@angular/router';
import { filter, Subscription } from 'rxjs';

@Component({
  selector: 'app-setup',
  standalone: true,
  imports: [CommonModule, RouterLink, RouterOutlet],
  template: `
    <div class="min-h-screen flex flex-col bg-gray-50 dark:bg-gray-900">
      <!-- Header -->
      <header class="sticky top-0 z-10 bg-white dark:bg-gray-800 shadow-sm py-4 px-4 sm:px-6 lg:px-8">
        <h1 class="text-center text-2xl font-bold text-gray-900 dark:text-white">Lethean VPN Setup</h1>
      </header>

      <!-- Main Content Area -->
      <main class="flex-grow flex items-center justify-center py-8 px-4 sm:px-6 lg:px-8">
        <div class="max-w-md w-full space-y-8">
          <ng-container *ngIf="!hasChildRoute">
            <!-- Initial setup options (buttons) -->
            <div>
              <h2 class="mt-6 text-center text-3xl font-extrabold text-gray-900 dark:text-white">
                Welcome to Lethean Setup
              </h2>
              <p class="mt-2 text-center text-sm text-gray-600 dark:text-gray-400">
                Please choose your setup option
              </p>
            </div>
            <div class="mt-8 space-y-4">
              <button type="button" routerLink="full" class="group relative w-full flex justify-center py-2 px-4 border border-transparent text-sm font-medium rounded-md text-white bg-indigo-600 hover:bg-indigo-700 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-indigo-500">
                Full Install
              </button>
              <button type="button" disabled class="group relative w-full flex justify-center py-2 px-4 border border-transparent text-sm font-medium rounded-md text-white bg-indigo-400 cursor-not-allowed opacity-50">
                Blockchain Only
              </button>
              <button type="button" disabled class="group relative w-full flex justify-center py-2 px-4 border border-transparent text-sm font-medium rounded-md text-white bg-indigo-400 cursor-not-allowed opacity-50">
                Gateway Client Only
              </button>
              <button type="button" disabled class="group relative w-full flex justify-center py-2 px-4 border border-transparent text-sm font-medium rounded-md text-white bg-indigo-400 cursor-not-allowed opacity-50">
                Seed Node
              </button>
            </div>
          </ng-container>
          <router-outlet *ngIf="hasChildRoute"></router-outlet>
        </div>
      </main>

      <!-- Footer with progress bar -->
      <footer class="sticky bottom-0 z-10 bg-white dark:bg-gray-800 shadow-sm py-4 px-4 sm:px-6 lg:px-8">
        <div>
          <h4 class="sr-only">Status</h4>
          <p class="text-sm font-medium text-gray-900 dark:text-white">Installing Lethean Desktop...</p>
          <div aria-hidden="true" class="mt-6">
            <div class="overflow-hidden rounded-full bg-gray-200 dark:bg-white/10">
              <div [style.width]="progressBarWidth" class="h-2 rounded-full bg-indigo-600 dark:bg-indigo-500"></div>
            </div>
            <div class="mt-6 hidden grid-cols-4 text-sm font-medium text-gray-600 sm:grid dark:text-gray-400">
              <a routerLink="." class="cursor-pointer" [class.text-indigo-600]="currentStep === 1" [class.dark:text-indigo-400]="currentStep === 1">Start Install</a>
              <a routerLink="full" class="text-center cursor-pointer" [class.text-indigo-600]="currentStep === 2" [class.dark:text-indigo-400]="currentStep === 2">Config</a>
              <div class="text-center" [class.text-indigo-600]="currentStep === 3" [class.dark:text-indigo-400]="currentStep === 3">Installing</div>
              <div class="text-right" [class.text-indigo-600]="currentStep === 4" [class.dark:text-indigo-400]="currentStep === 4">Welcome</div>
            </div>
          </div>
        </div>
      </footer>
    </div>
  `,
})
export class SetupComponent implements OnInit, OnDestroy {
  hasChildRoute: boolean = false;
  currentStep: number = 1; // 1: Start Install, 2: Config, 3: Installing, 4: Welcome
  progressBarWidth: string = '8%'; // Initial width for step 1
  private routerSubscription: Subscription | undefined;

  constructor(private router: Router, private activatedRoute: ActivatedRoute) {}

  ngOnInit(): void {
    this.checkChildRoute();
    this.routerSubscription = this.router.events.pipe(
      filter(event => event instanceof NavigationEnd)
    ).subscribe(() => {
      this.checkChildRoute();
      this.updateProgressBar();
    });
  }

  ngOnDestroy(): void {
    this.routerSubscription?.unsubscribe();
  }

  checkChildRoute(): void {
    this.hasChildRoute = this.activatedRoute.firstChild !== null;
  }

  updateProgressBar(): void {
    const currentPath = this.router.url;
    if (currentPath.includes('/setup/full') || currentPath.includes('/setup/blockchain') || currentPath.includes('/setup/gateway-client') || currentPath.includes('/setup/seed-node')) {
      this.currentStep = 2; // Config stage
      this.progressBarWidth = '33%';
    } else if (currentPath === '/setup') {
      this.currentStep = 1; // Start Install stage
      this.progressBarWidth = '8%';
    } else {
      // Default or other stages, can be expanded later
      this.currentStep = 1; // Fallback
      this.progressBarWidth = '8%';
    }
  }
}
