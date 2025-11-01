import {Component, CUSTOM_ELEMENTS_SCHEMA, OnDestroy, OnInit} from '@angular/core';
import { CommonModule, TitleCasePipe } from '@angular/common';
import { NavigationEnd, Router, RouterLink, RouterLinkActive, RouterOutlet } from '@angular/router';
import { ShowEnvironmentDialog } from "@lthn/core/display/service"
import { OpenDocsWindow } from "@lthn/docs/service"
import { EnableFeature, IsFeatureEnabled } from "@lthn/core/config/service";
import { TranslationService } from '../app/services/translation.service';
import { I18nService } from '../app/services/i18n.service';
import { Subscription } from 'rxjs';

@Component({
  selector: 'application-frame',
  standalone: true,
  schemas: [CUSTOM_ELEMENTS_SCHEMA],
  imports: [CommonModule, RouterOutlet, RouterLink, RouterLinkActive, TitleCasePipe],
  templateUrl: './application.frame.html',
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
