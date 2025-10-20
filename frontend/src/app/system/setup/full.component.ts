import {Component} from '@angular/core';
import {CommonModule} from '@angular/common';
import {FormsModule} from '@angular/forms';
import {Router} from '@angular/router';
import {SelectDirectory} from '@lthn/display/service';

@Component({
  selector: 'app-full-setup',
  standalone: true,
  imports: [CommonModule, FormsModule],
  template: `
    <div class="w-full">
      <h2 class="text-center text-3xl font-extrabold text-gray-900 dark:text-white mb-8">
        Full Installation Wizard
      </h2>

      <div class="lg:grid lg:grid-cols-3 lg:gap-8">
        <div class="lg:col-span-1">
          <nav aria-label="Progress">
            <ol role="list" class="overflow-hidden">
              <li class="relative pb-10">
                <div *ngIf="currentStep < 6" aria-hidden="true" class="absolute top-4 left-4 mt-0.5 -ml-px h-full w-0.5"
                     [ngClass]="{'bg-indigo-600 dark:bg-indigo-500': currentStep > 1, 'bg-gray-300 dark:bg-gray-700': currentStep <= 1}"></div>
                <!-- Step 1: Username -->
                <div class="group relative flex items-start">
                  <span class="flex h-9 items-center">
                    <span class="relative z-10 flex size-8 items-center justify-center rounded-full"
                          [ngClass]="{
                            'bg-indigo-600 group-hover:bg-indigo-800 dark:bg-indigo-500 dark:group-hover:bg-indigo-600': currentStep > 1,
                            'border-2 border-indigo-600 bg-white dark:border-indigo-500 dark:bg-gray-900': currentStep === 1,
                            'border-2 border-gray-300 bg-white group-hover:border-gray-400 dark:border-white/15 dark:bg-gray-900 dark:group-hover:border-white/25': currentStep < 1
                          }">
                      <svg *ngIf="currentStep > 1" viewBox="0 0 20 20" fill="currentColor" data-slot="icon"
                           aria-hidden="true" class="size-5 text-white">
                        <path
                          d="M16.704 4.153a.75.75 0 0 1 .143 1.052l-8 10.5a.75.75 0 0 1-1.127.075l-4.5-4.5a.75.75 0 0 1 1.06-1.06l3.894 3.893 7.48-9.817a.75.75 0 0 1 1.05-.143Z"
                          clip-rule="evenodd" fill-rule="evenodd"/>
                      </svg>
                      <span *ngIf="currentStep === 1"
                            class="size-2.5 rounded-full bg-indigo-600 dark:bg-indigo-500"></span>
                      <span *ngIf="currentStep < 1"
                            class="size-2.5 rounded-full bg-transparent group-hover:bg-gray-300 dark:group-hover:bg-white/15"></span>
                    </span>
                  </span>
                  <span class="ml-4 flex min-w-0 flex-col">
                    <div [ngClass]="{
                      'text-indigo-600 dark:text-indigo-400': currentStep === 1,
                      'text-gray-900 dark:text-white': currentStep > 1,
                      'text-gray-500 dark:text-gray-400': currentStep < 1
                    }" class="mt-4 flex rounded-md shadow-sm">
              <div class="relative flex flex-grow items-stretch focus-within:z-10">
                <input type="text" name="username" id="username" [(ngModel)]="username"
                       class="block w-full rounded-none rounded-l-md border-0 py-1.5 text-gray-900 ring-1 ring-inset ring-gray-300 placeholder:text-gray-400 focus:ring-2 focus:ring-inset focus:ring-indigo-600 sm:text-sm sm:leading-6 dark:bg-gray-700 dark:text-white dark:ring-gray-600 dark:placeholder:text-gray-400 dark:focus:ring-indigo-500"
                       placeholder=" Enter your username" [disabled]="currentStep > 1">
              </div>
              <button type="button" (click)="nextStep()" [disabled]="!username || currentStep > 1"
                      [ngClass]="{'text-white bg-indigo-600 hover:bg-indigo-700 dark:bg-indigo-500 dark:hover:bg-indigo-600': !(!username || currentStep > 1), 'text-gray-700 bg-gray-400 cursor-not-allowed': (!username || currentStep > 1)}"
                      class="relative -ml-px inline-flex items-center gap-x-1.5 rounded-r-md px-3 py-2 text-sm font-medium focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-indigo-500 dark:focus:ring-indigo-400">
                Next
              </button>
            </div>
                  </span>
                </div>
              </li>

              <li class="relative pb-10">
                <div *ngIf="currentStep < 6" aria-hidden="true" class="absolute top-4 left-4 mt-0.5 -ml-px h-full w-0.5"
                     [ngClass]="{'bg-indigo-600 dark:bg-indigo-500': currentStep > 2, 'bg-gray-300 dark:bg-gray-700': currentStep <= 2}"></div>
                <!-- Step 2: Install Directory -->
                <div class="group relative flex items-start">
                  <span class="flex h-9 items-center">
                    <span class="relative z-10 flex size-8 items-center justify-center rounded-full"
                          [ngClass]="{
                            'bg-indigo-600 group-hover:bg-indigo-800 dark:bg-indigo-500 dark:group-hover:bg-indigo-600': currentStep > 2,
                            'border-2 border-indigo-600 bg-white dark:border-indigo-500 dark:bg-gray-900': currentStep === 2,
                            'border-2 border-gray-300 bg-white group-hover:border-gray-400 dark:border-white/15 dark:bg-gray-900 dark:group-hover:border-white/25': currentStep < 2
                          }">
                      <svg *ngIf="currentStep > 2" viewBox="0 0 20 20" fill="currentColor" data-slot="icon"
                           aria-hidden="true" class="size-5 text-white">
                        <path
                          d="M16.704 4.153a.75.75 0 0 1 .143 1.052l-8 10.5a.75.75 0 0 1-1.127.075l-4.5-4.5a.75.75 0 0 1 1.06-1.06l3.894 3.893 7.48-9.817a.75.75 0 0 1 1.05-.143Z"
                          clip-rule="evenodd" fill-rule="evenodd"/>
                      </svg>
                      <span *ngIf="currentStep === 2"
                            class="size-2.5 rounded-full bg-indigo-600 dark:bg-indigo-500"></span>
                      <span *ngIf="currentStep < 2"
                            class="size-2.5 rounded-full bg-transparent group-hover:bg-gray-300 dark:group-hover:bg-white/15"></span>
                    </span>
                  </span>
                  <span class="ml-4 flex min-w-0 flex-col">
                    <div [ngClass]="{
                      'text-indigo-600 dark:text-indigo-400': currentStep === 2,
                      'text-gray-900 dark:text-white': currentStep > 2,
                      'text-gray-500 dark:text-gray-400': currentStep < 2
                    }" class="mt-4 flex rounded-md shadow-sm">
              <div class="relative flex flex-grow items-stretch focus-within:z-10">
                <input type="text" name="installDirectory" id="installDirectory" [(ngModel)]="installDirectory"
                       class="block w-full rounded-none rounded-l-md border-0 py-1.5 text-gray-900 ring-1 ring-inset ring-gray-300 placeholder:text-gray-400 focus:ring-2 focus:ring-inset focus:ring-indigo-600 sm:text-sm sm:leading-6 dark:bg-gray-700 dark:text-white dark:ring-gray-600 dark:placeholder:text-gray-400 dark:focus:ring-indigo-500"
                       placeholder=" e.g., ~/lethean" [disabled]="currentStep > 2">
                <button type="button" (click)="selectInstallDirectory()" [disabled]="currentStep > 2"
                        class="relative -ml-px inline-flex items-center gap-x-1.5 rounded-r-md px-3 py-2 text-sm font-medium text-white bg-gray-600 hover:bg-gray-700 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-gray-500 dark:bg-gray-500 dark:hover:bg-gray-600 dark:focus:ring-gray-400">
                  Browse
                </button>
              </div>
              <button type="button" (click)="nextStep()" [disabled]="!installDirectory || currentStep > 2"
                      [ngClass]="{'text-white bg-indigo-600 hover:bg-indigo-700 dark:bg-indigo-500 dark:hover:bg-indigo-600': !(!installDirectory || currentStep > 2), 'text-gray-700 bg-gray-400 cursor-not-allowed': (!installDirectory || currentStep > 2)}"
                      class="relative -ml-px inline-flex items-center gap-x-1.5 rounded-r-md px-3 py-2 text-sm font-medium focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-indigo-500 dark:focus:ring-indigo-400">
                Next
              </button>
            </div>
                  </span>
                </div>
              </li>

              <li class="relative pb-10">
                <div *ngIf="currentStep < 6" aria-hidden="true" class="absolute top-4 left-4 mt-0.5 -ml-px h-full w-0.5"
                     [ngClass]="{'bg-indigo-600 dark:bg-indigo-500': currentStep > 3, 'bg-gray-300 dark:bg-gray-700': currentStep <= 3}"></div>
                <!-- Step 3: Blockchain -->
                <div class="group relative flex items-start">
                  <span class="flex h-9 items-center">
                    <span class="relative z-10 flex size-8 items-center justify-center rounded-full"
                          [ngClass]="{
                            'bg-indigo-600 group-hover:bg-indigo-800 dark:bg-indigo-500 dark:group-hover:bg-indigo-600': currentStep > 3,
                            'border-2 border-indigo-600 bg-white dark:border-indigo-500 dark:bg-gray-900': currentStep === 3,
                            'border-2 border-gray-300 bg-white group-hover:border-gray-400 dark:border-white/15 dark:bg-gray-900 dark:group-hover:border-white/25': currentStep < 3
                          }">
                      <svg *ngIf="currentStep > 3" viewBox="0 0 20 20" fill="currentColor" data-slot="icon"
                           aria-hidden="true" class="size-5 text-white">
                        <path
                          d="M16.704 4.153a.75.75 0 0 1 .143 1.052l-8 10.5a.75.75 0 0 1-1.127.075l-4.5-4.5a.75.75 0 0 1 1.06-1.06l3.894 3.893 7.48-9.817a.75.75 0 0 1 1.05-.143Z"
                          clip-rule="evenodd" fill-rule="evenodd"/>
                      </svg>
                      <span *ngIf="currentStep === 3"
                            class="size-2.5 rounded-full bg-indigo-600 dark:bg-indigo-500"></span>
                      <span *ngIf="currentStep < 3"
                            class="size-2.5 rounded-full bg-transparent group-hover:bg-gray-300 dark:group-hover:bg-white/15"></span>
                    </span>
                  </span>
                  <span class="ml-4 flex min-w-0 flex-col">
                    <span class="text-sm font-medium" [ngClass]="{
                      'text-indigo-600 dark:text-indigo-400': currentStep === 3,
                      'text-gray-900 dark:text-white': currentStep > 3,
                      'text-gray-500 dark:text-gray-400': currentStep < 3
                    }">Blockchain</span>
                    <span class="text-sm text-gray-500 dark:text-gray-400">Configure blockchain settings.</span>
                  </span>
                </div>
              </li>

              <li class="relative pb-10">
                <div *ngIf="currentStep < 6" aria-hidden="true" class="absolute top-4 left-4 mt-0.5 -ml-px h-full w-0.5"
                     [ngClass]="{'bg-indigo-600 dark:bg-indigo-500': currentStep > 4, 'bg-gray-300 dark:bg-gray-700': currentStep <= 4}"></div>
                <!-- Step 4: Gateway Client -->
                <div class="group relative flex items-start">
                  <span class="flex h-9 items-center">
                    <span class="relative z-10 flex size-8 items-center justify-center rounded-full"
                          [ngClass]="{
                            'bg-indigo-600 group-hover:bg-indigo-800 dark:bg-indigo-500 dark:group-hover:bg-indigo-600': currentStep > 4,
                            'border-2 border-indigo-600 bg-white dark:border-indigo-500 dark:bg-gray-900': currentStep === 4,
                            'border-2 border-gray-300 bg-white group-hover:border-gray-400 dark:border-white/15 dark:bg-gray-900 dark:group-hover:border-white/25': currentStep < 4
                          }">
                      <svg *ngIf="currentStep > 4" viewBox="0 0 20 20" fill="currentColor" data-slot="icon"
                           aria-hidden="true" class="size-5 text-white">
                        <path
                          d="M16.704 4.153a.75.75 0 0 1 .143 1.052l-8 10.5a.75.75 0 0 1-1.127.075l-4.5-4.5a.75.75 0 0 1 1.06-1.06l3.894 3.893 7.48-9.817a.75.75 0 0 1 1.05-.143Z"
                          clip-rule="evenodd" fill-rule="evenodd"/>
                      </svg>
                      <span *ngIf="currentStep === 4"
                            class="size-2.5 rounded-full bg-indigo-600 dark:bg-indigo-500"></span>
                      <span *ngIf="currentStep < 4"
                            class="size-2.5 rounded-full bg-transparent group-hover:bg-gray-300 dark:group-hover:bg-white/15"></span>
                    </span>
                  </span>
                  <span class="ml-4 flex min-w-0 flex-col">
                    <span class="text-sm font-medium" [ngClass]="{
                      'text-indigo-600 dark:text-indigo-400': currentStep === 4,
                      'text-gray-900 dark:text-white': currentStep > 4,
                      'text-gray-500 dark:text-gray-400': currentStep < 4
                    }">Gateway Client</span>
                    <span class="text-sm text-gray-500 dark:text-gray-400">Setup gateway client.</span>
                  </span>
                </div>
              </li>

              <li class="relative pb-10">
                <div *ngIf="currentStep < 6" aria-hidden="true" class="absolute top-4 left-4 mt-0.5 -ml-px h-full w-0.5"
                     [ngClass]="{'bg-indigo-600 dark:bg-indigo-500': currentStep > 5, 'bg-gray-300 dark:bg-gray-700': currentStep <= 5}"></div>
                <!-- Step 5: Seed Node -->
                <div class="group relative flex items-start">
                  <span class="flex h-9 items-center">
                    <span class="relative z-10 flex size-8 items-center justify-center rounded-full"
                          [ngClass]="{
                            'bg-indigo-600 group-hover:bg-indigo-800 dark:bg-indigo-500 dark:group-hover:bg-indigo-600': currentStep > 5,
                            'border-2 border-indigo-600 bg-white dark:border-indigo-500 dark:bg-gray-900': currentStep === 5,
                            'border-2 border-gray-300 bg-white group-hover:border-gray-400 dark:border-white/15 dark:bg-gray-900 dark:group-hover:border-white/25': currentStep < 5
                          }">
                      <svg *ngIf="currentStep > 5" viewBox="0 0 20 20" fill="currentColor" data-slot="icon"
                           aria-hidden="true" class="size-5 text-white">
                        <path
                          d="M16.704 4.153a.75.75 0 0 1 .143 1.052l-8 10.5a.75.75 0 0 1-1.127.075l-4.5-4.5a.75.75 0 0 1 1.06-1.06l3.894 3.893 7.48-9.817a.75.75 0 0 1 1.05-.143Z"
                          clip-rule="evenodd" fill-rule="evenodd"/>
                      </svg>
                      <span *ngIf="currentStep === 5"
                            class="size-2.5 rounded-full bg-indigo-600 dark:bg-indigo-500"></span>
                      <span *ngIf="currentStep < 5"
                            class="size-2.5 rounded-full bg-transparent group-hover:bg-gray-300 dark:group-hover:bg-white/15"></span>
                    </span>
                  </span>
                  <span class="ml-4 flex min-w-0 flex-col">
                    <span class="text-sm font-medium" [ngClass]="{
                      'text-indigo-600 dark:text-indigo-400': currentStep === 5,
                      'text-gray-900 dark:text-white': currentStep > 5,
                      'text-gray-500 dark:text-gray-400': currentStep < 5
                    }">Seed Node</span>
                    <span class="text-sm text-gray-500 dark:text-gray-400">Configure seed node.</span>
                  </span>
                </div>
              </li>

              <li class="relative">
                <!-- Step 6: Complete -->
                <div class="group relative flex items-start">
                  <span class="flex h-9 items-center">
                    <span class="relative z-10 flex size-8 items-center justify-center rounded-full"
                          [ngClass]="{
                            'bg-indigo-600 group-hover:bg-indigo-800 dark:bg-indigo-500 dark:group-hover:bg-indigo-600': currentStep > 6,
                            'border-2 border-indigo-600 bg-white dark:border-indigo-500 dark:bg-gray-900': currentStep === 6,
                            'border-2 border-gray-300 bg-white group-hover:border-gray-400 dark:border-white/15 dark:bg-gray-900 dark:group-hover:border-white/25': currentStep < 6
                          }">
                      <svg *ngIf="currentStep > 6" viewBox="0 0 20 20" fill="currentColor" data-slot="icon"
                           aria-hidden="true" class="size-5 text-white">
                        <path
                          d="M16.704 4.153a.75.75 0 0 1 .143 1.052l-8 10.5a.75.75 0 0 1-1.127.075l-4.5-4.5a.75.75 0 0 1 1.06-1.06l3.894 3.893 7.48-9.817a.75.75 0 0 1 1.05-.143Z"
                          clip-rule="evenodd" fill-rule="evenodd"/>
                      </svg>
                      <span *ngIf="currentStep === 6"
                            class="size-2.5 rounded-full bg-indigo-600 dark:bg-indigo-500"></span>
                      <span *ngIf="currentStep < 6"
                            class="size-2.5 rounded-full bg-transparent group-hover:bg-gray-300 dark:group-hover:bg-white/15"></span>
                    </span>
                  </span>
                  <span class="ml-4 flex min-w-0 flex-col">
                    <span class="text-sm font-medium" [ngClass]="{
                      'text-indigo-600 dark:text-indigo-400': currentStep === 6,
                      'text-gray-900 dark:text-white': currentStep > 6,
                      'text-gray-500 dark:text-gray-400': currentStep < 6
                    }">Complete</span>
                    <span class="text-sm text-gray-500 dark:text-gray-400">Finish setup.</span>
                  </span>
                </div>
              </li>
            </ol>
          </nav>
        </div>

        <div class="mt-10 lg:col-span-2 lg:mt-0">
          <div *ngIf="currentStep === 3" class="space-y-6">
            <h3 class="text-lg font-medium leading-6 text-gray-900 dark:text-white">Blockchain Configuration</h3>
            <p class="mt-1 text-sm text-gray-500 dark:text-gray-400">Details for blockchain setup...</p>
            <button type="button" (click)="nextStep()" [disabled]="currentStep > 3"
                    [ngClass]="{'text-white bg-indigo-600 hover:bg-indigo-700 dark:bg-indigo-500 dark:hover:bg-indigo-600': !(currentStep > 3), 'text-gray-700 bg-gray-400 cursor-not-allowed': (currentStep > 3)}"
                    class="inline-flex items-center px-4 py-2 border border-transparent text-sm font-medium rounded-md shadow-sm focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-indigo-500 dark:focus:ring-indigo-400">
              Next
            </button>
          </div>

          <div *ngIf="currentStep === 4" class="space-y-6">
            <h3 class="text-lg font-medium leading-6 text-gray-900 dark:text-white">Gateway Client Setup</h3>
            <p class="mt-1 text-sm text-gray-500 dark:text-gray-400">Details for gateway client setup...</p>
            <button type="button" (click)="nextStep()" [disabled]="currentStep > 4"
                    [ngClass]="{'text-white bg-indigo-600 hover:bg-indigo-700 dark:bg-indigo-500 dark:hover:bg-indigo-600': !(currentStep > 4), 'text-gray-700 bg-gray-400 cursor-not-allowed': (currentStep > 4)}"
                    class="inline-flex items-center px-4 py-2 border border-transparent text-sm font-medium rounded-md shadow-sm focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-indigo-500 dark:focus:ring-indigo-400">
              Next
            </button>
          </div>

          <div *ngIf="currentStep === 5" class="space-y-6">
            <h3 class="text-lg font-medium leading-6 text-gray-900 dark:text-white">Seed Node Configuration</h3>
            <p class="mt-1 text-sm text-gray-500 dark:text-gray-400">Details for seed node setup...</p>
            <button type="button" (click)="nextStep()" [disabled]="currentStep > 5"
                    [ngClass]="{'text-white bg-indigo-600 hover:bg-indigo-700 dark:bg-indigo-500 dark:hover:bg-indigo-600': !(currentStep > 5), 'text-gray-700 bg-gray-400 cursor-not-allowed': (currentStep > 5)}"
                    class="inline-flex items-center px-4 py-2 border border-transparent text-sm font-medium rounded-md shadow-sm focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-indigo-500 dark:focus:ring-indigo-400">
              Next
            </button>
          </div>

          <div *ngIf="currentStep === 6" class="space-y-6">
            <h3 class="text-lg font-medium leading-6 text-gray-900 dark:text-white">Setup Complete!</h3>
            <p class="mt-1 text-sm text-gray-500 dark:text-gray-400">You have successfully completed the setup
              wizard.</p>
            <button type="button" (click)="finishSetup()"
                    class="inline-flex items-center px-4 py-2 border border-transparent text-sm font-medium rounded-md shadow-sm text-white bg-green-600 hover:bg-green-700 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-green-500 dark:bg-green-500 dark:hover:bg-green-600">
              Go to Dashboard
            </button>
          </div>
        </div>
      </div>
    </div>
  `,
})
export class FullComponent {
  currentStep: number = 1;
  username: string = '';
  installDirectory: string = '~/lethean'; // New property for install directory

  constructor(private router: Router) {
  }

  async selectInstallDirectory(): Promise<void> {
    try {
      const selectedPath = await SelectDirectory();
      if (selectedPath) {
        this.installDirectory = selectedPath;
      }
    } catch (error) {
      console.error('Error selecting directory:', error);
    }
  }

  nextStep(): void {
    if (this.currentStep < 6) { // Updated total steps
      this.currentStep++;
    }
  }

  finishSetup(): void {
    console.log('Setup finished!');
    console.log('Username:', this.username);
    console.log('Install Directory:', this.installDirectory);
    this.router.navigate(['/blockchain']);
  }
}
