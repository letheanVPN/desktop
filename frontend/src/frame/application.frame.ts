import { Component, OnDestroy, OnInit } from '@angular/core';
import { CommonModule, TitleCasePipe } from '@angular/common';
import { NavigationEnd, Router, RouterLink, RouterLinkActive, RouterOutlet } from '@angular/router';
import { ShowEnvironmentDialog } from "@lthn/display/service"
import { OpenDocsWindow } from "@lthn/docs/service"
import { EnableFeature, IsFeatureEnabled } from "@lthn/config/service";
import { TranslationService } from '../app/services/translation.service';
import { I18nService } from '../app/services/i18n.service';
import { Subscription } from 'rxjs';

@Component({
  selector: 'application-frame',
  standalone: true,
  imports: [CommonModule, RouterOutlet, RouterLink, RouterLinkActive, TitleCasePipe],
  template: `
    <!-- Mobile sidebar -->
    @if (sidebarOpen) {
      <div class="relative z-50 lg:hidden" role="dialog" aria-modal="true">
        <div class="fixed inset-0 bg-gray-900/80"></div>

        <div class="fixed inset-0 flex">
          <div class="relative mr-16 flex w-full max-w-xs flex-1">
            <div class="absolute top-0 left-full flex w-16 justify-center pt-5">
              <button type="button" (click)="sidebarOpen = false" class="-m-2.5 p-2.5">
                <span class="sr-only">Close sidebar</span>
                <i class="fa-regular fa-xmark fa-2xl text-white"></i>
              </button>
            </div>
            <div class="flex grow flex-col gap-y-5 overflow-y-auto bg-gray-900 px-6 pb-2 ring-1 ring-white/10">
              <div class="flex h-16 shrink-0 items-center">
                <img class="h-58 w-80 pt-4" src="./logo/lthn/logo-full-gradient.png" alt="Lethean Community">
              </div>
              <nav class="flex flex-1 flex-col">
                <ul role="list" class="-mx-2 flex-1 space-y-1">
                  @for (item of navigation; track item.name) {
                    <li>
                      <a [routerLink]="item.href" routerLinkActive="bg-gray-800 text-white"
                         [routerLinkActiveOptions]="{exact: true}"
                         class="text-gray-400 hover:bg-gray-800 hover:text-white group flex justify-center items-center gap-x-3 rounded-md p-4 text-sm/6 font-semibold">
                        <i [class]="item.icon"></i>
                        {{ item.name }}
                      </a>
                    </li>
                  }
                </ul>
              </nav>
            </div>
          </div>
        </div>
      </div>
    }

    <!-- Static sidebar for desktop -->
    <div
      class="hidden lg:fixed lg:inset-y-0 lg:left-0 lg:z-50 lg:block lg:w-20 lg:overflow-y-auto lg:bg-gray-900 lg:pb-4 dark:before:pointer-events-none dark:before:absolute dark:before:inset-0 dark:before:border-r dark:before:border-white/10 dark:before:bg-black/10">
      <div class="relative flex h-16 shrink-0 items-center justify-center">
        <img class="mt-1 h-20 w-auto" src="./logo/lthn/logo-icon-gradient.png" alt="Lethean Community">
      </div>
      <nav class="relative mt-1">
        <ul role="list" class="flex flex-col items-center space-y-1">
          @for (item of this.navigation; track item.name) {
            <li>
              <a [routerLink]="item.href" routerLinkActive="bg-white/5 text-white" [routerLinkActiveOptions]="{exact: true}"
                 class="text-gray-400 hover:bg-white/5 hover:text-white group flex justify-center items-center rounded-md p-4 text-sm/6 font-semibold h-16">
                <i [class]="item.icon"></i>
                <span class="sr-only">{{ item.name }}</span>
              </a>
            </li>
          }
        </ul>
      </nav>
    </div>

    <div class="lg:pl-20 pb-10">
      <div
        class="sticky top-0 z-40 flex h-16 shrink-0 items-center gap-x-4 border-b border-gray-200 bg-white px-4 shadow-xs sm:gap-x-6 sm:px-6 lg:px-8 dark:border-white/10 dark:bg-gray-900 dark:shadow-none dark:before:pointer-events-none dark:before:absolute dark:before:inset-0 dark:before:bg-black/10">
        <button type="button" (click)="sidebarOpen = true"
                class="-m-2.5 p-2.5 text-gray-700 lg:hidden dark:text-gray-400">
          <span class="sr-only">Open sidebar</span>
          <i class="fa-regular fa-bars fa-2xl"></i>
        </button>

        <!-- Separator -->
        <div aria-hidden="true" class="h-6 w-px bg-gray-900/10 lg:hidden dark:bg-white/10"></div>

        <div class="flex flex-1 gap-x-4 self-stretch lg:gap-x-6">
          <form action="#" method="GET" class="grid flex-1 grid-cols-1">
            <input name="search" placeholder="{{ t._('app.core.ui.search') }}" aria-label="Search"
                   class="col-start-1 row-start-1 block size-full bg-white pl-8 text-base text-gray-900 outline-hidden placeholder:text-gray-400 sm:text-sm/6 dark:bg-gray-900 dark:text-white dark:placeholder:text-gray-500"/>
            <i class="fa-light fa-xl fa-magnifying-glass pointer-events-none col-start-1 row-start-1 self-center text-gray-400"></i>
          </form>
          <div class="flex items-center gap-x-4 lg:gap-x-6">
            <button type="button" class="-m-2.5 p-2.5 text-gray-400 hover:text-gray-500 dark:hover:text-white">
              <span class="sr-only">View notifications</span>
              <i class="fa-light fa-bell fa-xl"></i>
            </button>

            <!-- Separator -->
            <div aria-hidden="true" class="hidden lg:block lg:h-6 lg:w-px bg-gray-900/10 dark:lg:bg-white/10"></div>

            <!-- Profile dropdown -->
            <div class="relative">
              <button (click)="userMenuOpen = !userMenuOpen" class="relative flex items-center">
                <span class="absolute -inset-1.5"></span>
                <span class="sr-only">Open user menu</span>
                <img src="avatar.png" alt="Current User Avatar Image" class="size-8 rounded-full bg-gray-50 outline -outline-offset-1 outline-black/5 dark:bg-gray-800 dark:outline-white/10"/>
                <span class="hidden lg:flex lg:items-center">
                  <span aria-hidden="true" class="ml-4 w-32 text-sm/6 font-semibold text-gray-900 dark:text-white">{{ currentRole }} Hub</span>
                  <i class="ml-2 fa-regular fa-chevron-down text-sm/6 text-gray-400"></i>
                </span>
              </button>
              @if (userMenuOpen) {
                <div
                  class="absolute right-0 z-10 mt-2.5 w-48 origin-top-right rounded-md bg-white py-2 shadow-lg ring-1 ring-black ring-opacity-5 focus:outline-none dark:bg-gray-800 dark:ring-white/10">
                  <a (click)="openDocs()"
                     class="flex items-center gap-x-3 px-3 py-1 text-sm/6 text-gray-900 focus:bg-gray-50 focus:outline-hidden dark:text-white dark:focus:bg-white/5">
                    <i class="fa-regular fa-code"></i>
                    Documentation
                  </a>
                  @for (item of userNavigation; track item.name) {
                    <a [routerLink]="item.href" (click)="userMenuOpen = false" class="flex items-center gap-x-3 px-3 py-1 text-sm/6 text-gray-900 focus:bg-gray-50 focus:outline-hidden dark:text-white dark:focus:bg-white/5 h-8 ">
                      <i [class]="item.icon"></i>
                      {{ item.name }}
                    </a>
                  }
                  <div class="my-2 h-px bg-gray-200 dark:bg-white/10"></div>
                  <div
                    class="px-3 py-2 text-xs font-semibold text-gray-500 uppercase tracking-wider dark:text-gray-400">
                    Switch Role
                  </div>
                  @for (item of roleNavigation; track item.name) {
                    <a [routerLink]="item.href" (click)="switchRole(item.name)"
                       class="block px-3 py-1 text-sm/6 text-gray-900 focus:bg-gray-50 focus:outline-hidden dark:text-white dark:focus:bg-white/5">
                      {{ item.name }}
                    </a>
                  }
                </div>
              }
            </div>
          </div>
        </div>
      </div>

      <main >
        <div class="px-0 py-0">
          @if (featureKey && !isFeatureEnabled) {
            <div class="bg-yellow-100 border-l-4 border-yellow-500 text-yellow-700 p-4" role="alert">
              <p class="font-bold">{{ featureKey | titlecase }} is Disabled</p>
              <p>Activate {{ featureKey }} to access its features.</p>
              <button (click)="activateFeature()"
                      class="mt-4 px-4 py-2 bg-blue-500 text-white rounded hover:bg-blue-600">
                Activate {{ featureKey | titlecase }}
              </button>
            </div>
          } @else {
            <router-outlet></router-outlet>
          }
        </div>
      </main>
    </div>


    <footer
      class="fixed inset-x-0 bottom-0 z-40 h-10 border-t border-gray-200 bg-white dark:border-white/10 dark:bg-gray-900 lg:left-20">
      <div class="flex h-full items-center justify-between px-4 sm:px-6 lg:px-8">
        <button class="text-sm text-gray-500 dark:text-gray-400 cursor-pointer" (click)="ShowEnvironmentDialog()">
          v0.0.1 <i class="fa-brands fa-pied-piper-alt fa-xl"></i>
        </button>
        @if (currentRole === 'Developer') {
          <button class="ml-4 text-sm text-gray-500 dark:text-gray-400 cursor-pointer" (click)="showTestDialog()">Test
            Dialog
          </button>
        }
        <span class="text-sm text-gray-500 dark:text-gray-400">{{ time }}</span>
      </div>
    </footer>
  `,
})
export class ApplicationFrame implements OnInit, OnDestroy {
  sidebarOpen = false;
  userMenuOpen = false;
  currentRole = 'Developer';
  time: string = '';
  private intervalId: number | undefined;
  private langChangeSubscription: Subscription | undefined;

  featureKey: string | null = null;
  isFeatureEnabled: boolean = false;
  userNavigation: any[] = [];
  navigation: any[] = [];
  roleNavigation: any[] = [];

  constructor(
    private router: Router,
    public t: TranslationService,
    private i18nService: I18nService
    ) {


  }

  async ngOnInit(): Promise<void> {
    this.updateTime();
    this.intervalId = window.setInterval(() => {
      this.updateTime();
    }, 1000);

    await this.t.onReady();
    this.initializeUserNavigation();

    this.langChangeSubscription = this.i18nService.currentLanguage$.subscribe(async () => {
        await this.t.onReady();
        this.initializeUserNavigation();
    });

    this.router.events.subscribe(event => {
      if (event instanceof NavigationEnd) {
        this.extractFeatureKeyAndCheckStatus(event.urlAfterRedirects);
      }
    });
    this.navigation = [
      { name: this.t._('menu.dashboard'), href: '', icon: "fa-grid-round fa-regular fa-2xl shrink-0" },
      { name: this.t._('menu.mining'), href: 'mining', icon: "fa-pickaxe fa-regular fa-2xl shrink-0" },
      { name: this.t._('menu.blockchain'), href: 'blockchain', icon: "fa-chart-network fa-regular fa-2xl shrink-0" },
      { name: this.t._('Developer'), href: 'dev/edit', icon: "fa-code fa-regular fa-2xl shrink-0" },
      { name: this.t._('Networking'), href: 'net', icon: "fa-network-wired fa-regular fa-2xl shrink-0" },
      { name: this.t._('Settings'), href: 'system/settings', icon: "fa-gear-code fa-regular fa-2xl shrink-0" },
    ];

    this.roleNavigation = [
      { name: this.t._('menu.hub-client'), href: '/config/client-hub' },
      { name: this.t._('menu.hub-server'), href: '/config/server-hub' },
      { name: this.t._('menu.hub-developer'), href: '/config/developer-hub' },
      { name: this.t._('menu.hub-gateway'), href: '/config/gateway-hub' },
      { name: this.t._('menu.hub-admin'), href: '/config/admin-hub' },
    ];
    await this.extractFeatureKeyAndCheckStatus(this.router.url); // Initial check
  }

  ngOnDestroy(): void {
    if (this.intervalId) {
      clearInterval(this.intervalId);
    }
    if (this.langChangeSubscription) {
        this.langChangeSubscription.unsubscribe();
    }
  }

  initializeUserNavigation(): void {
    this.userNavigation = [
      { name: this.t._('menu.your-profile'), href: '#', icon: "fa-id-card fa-regular" },
      { name: this.t._('menu.logout'), href: '#', icon: "fa-right-from-bracket fa-regular" },
    ];
  }

  updateTime(): void {
    const now = new Date();
    this.time = now.toLocaleTimeString();
  }

  async extractFeatureKeyAndCheckStatus(url: string): Promise<void> {
    // Remove leading slash and split by slash
    const parts = url.startsWith('/') ? url.substring(1).split('/') : url.split('/');
    if (parts.length > 0 && parts[0] !== '') {
      this.featureKey = parts[0];
      await this.checkFeatureStatus();
    } else {
      this.featureKey = null;
      this.isFeatureEnabled = true; // No feature key, so assume enabled
    }
  }

  async checkFeatureStatus(): Promise<void> {
    if (this.featureKey) {
      try {
        this.isFeatureEnabled = await IsFeatureEnabled(this.featureKey);
      } catch (error) {
        console.error(`Error checking feature ${this.featureKey}:`, error);
        this.isFeatureEnabled = false;
      }
    } else {
      this.isFeatureEnabled = true;
    }
  }

  async activateFeature(): Promise<void> {
    if (this.featureKey) {
      try {
        await EnableFeature(this.featureKey);
        await this.checkFeatureStatus();
      } catch (error) {
        console.error(`Error activating feature ${this.featureKey}:`, error);
      }
    }
  }

  showTestDialog(): void {
    alert('Test Dialog Triggered!');
  }

  openDocs() {
    return OpenDocsWindow("getting-started/chain#using-the-cli")
  }
  switchRole(roleName: string) {
    if (roleName.endsWith(' Hub')) {
      this.currentRole = roleName.replace(' Hub', '');
    }
    this.userMenuOpen = false;
  }

  protected readonly ShowEnvironmentDialog = ShowEnvironmentDialog;
}
