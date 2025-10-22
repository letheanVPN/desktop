import {Component, signal} from '@angular/core';
import {FormsModule} from '@angular/forms';
import {Router} from '@angular/router';
import {SelectDirectory} from '@lthn/core/display/service';
import {Save} from '@lthn/core/config/service';

@Component({
  selector: 'app-full-setup',
  standalone: true,
  imports: [FormsModule],
  template: `
    <div class="w-full">
      <h2 class="text-center text-3xl font-extrabold text-gray-900 dark:text-white mb-8">
        Full Installation Wizard
      </h2>

      <div class="lg:grid lg:grid-cols-3 lg:gap-8">
        <div class="lg:col-span-1">
          <nav aria-label="Progress">
            <ol role="list" class="overflow-hidden">
<!--              <li class="relative pb-10">-->
<!--                @if (currentStep() < 6) {-->
<!--                  <div aria-hidden="true" class="absolute top-4 left-4 mt-0.5 -ml-px h-full w-0.5"-->
<!--                       [class.bg-indigo-600]="currentStep() > 1"-->
<!--                       [class.dark:bg-indigo-500]="currentStep() > 1"-->
<!--                       [class.bg-gray-300]="currentStep() <= 1"-->
<!--                       [class.dark:bg-gray-700]="currentStep() <= 1"></div>-->
<!--                }-->
<!--                &lt;!&ndash; Step 1: Username &ndash;&gt;-->
<!--                <div class="group relative flex items-start">-->
<!--                  <span class="flex h-9 items-center">-->
<!--                    <span class="relative z-10 flex size-8 items-center justify-center rounded-full"-->
<!--                          [class.bg-indigo-600]="currentStep() > 1"-->
<!--                          [class.group-hover:bg-indigo-800]="currentStep() > 1"-->
<!--                          [class.dark:bg-indigo-500]="currentStep() > 1"-->
<!--                          [class.dark:group-hover:bg-indigo-600]="currentStep() > 1"-->
<!--                          [class.border-2]="currentStep() <= 1"-->
<!--                          [class.border-indigo-600]="currentStep() === 1"-->
<!--                          [class.bg-white]="currentStep() === 1"-->
<!--                          [class.dark:border-indigo-500]="currentStep() === 1"-->
<!--                          [class.dark:bg-gray-900]="currentStep() === 1"-->
<!--                          [class.border-gray-300]="currentStep() < 1"-->
<!--                          [class.group-hover:border-gray-400]="currentStep() < 1"-->
<!--                          [class.dark:border-white&#47;15]="currentStep() < 1"-->
<!--                          [class.dark:group-hover:border-white&#47;25]="currentStep() < 1">-->
<!--                      @if (currentStep() > 1) {-->
<!--                        <svg viewBox="0 0 20 20" fill="currentColor" data-slot="icon"-->
<!--                             aria-hidden="true" class="size-5 text-white">-->
<!--                          <path-->
<!--                            d="M16.704 4.153a.75.75 0 0 1 .143 1.052l-8 10.5a.75.75 0 0 1-1.127.075l-4.5-4.5a.75.75 0 0 1 1.06-1.06l3.894 3.893 7.48-9.817a.75.75 0 0 1 1.05-.143Z"-->
<!--                            clip-rule="evenodd" fill-rule="evenodd"/>-->
<!--                        </svg>-->
<!--                      }-->
<!--                      @if (currentStep() === 1) {-->
<!--                        <span class="size-2.5 rounded-full bg-indigo-600 dark:bg-indigo-500"></span>-->
<!--                      }-->
<!--                      @if (currentStep() < 1) {-->
<!--                        <span class="size-2.5 rounded-full bg-transparent group-hover:bg-gray-300 dark:group-hover:bg-white&#47;15"></span>-->
<!--                      }-->
<!--                    </span>-->
<!--                  </span>-->
<!--                  <span class="ml-4 flex min-w-0 flex-col">-->
<!--                    <div [class.text-indigo-600]="currentStep() === 1"-->
<!--                         [class.dark:text-indigo-400]="currentStep() === 1"-->
<!--                         [class.text-gray-900]="currentStep() > 1"-->
<!--                         [class.dark:text-white]="currentStep() > 1"-->
<!--                         [class.text-gray-500]="currentStep() < 1"-->
<!--                         [class.dark:text-gray-400]="currentStep() < 1"-->
<!--                         class="mt-4 flex rounded-md shadow-sm">-->
<!--              <div class="relative flex flex-grow items-stretch focus-within:z-10">-->
<!--                <input type="text" name="username" id="username" [(ngModel)]="username"-->
<!--                       class="block w-full rounded-none rounded-l-md border-0 py-1.5 text-gray-900 ring-1 ring-inset ring-gray-300 placeholder:text-gray-400 focus:ring-2 focus:ring-inset focus:ring-indigo-600 sm:text-sm sm:leading-6 dark:bg-gray-700 dark:text-white dark:ring-gray-600 dark:placeholder:text-gray-400 dark:focus:ring-indigo-500"-->
<!--                       placeholder=" Enter your username" [disabled]="currentStep() > 1">-->
<!--              </div>-->
<!--              <button type="button" (click)="nextStep()" [disabled]="!username || currentStep() > 1"-->
<!--                      [class.text-white]="username && currentStep() <= 1"-->
<!--                      [class.bg-indigo-600]="username && currentStep() <= 1"-->
<!--                      [class.hover:bg-indigo-700]="username && currentStep() <= 1"-->
<!--                      [class.dark:bg-indigo-500]="username && currentStep() <= 1"-->
<!--                      [class.dark:hover:bg-indigo-600]="username && currentStep() <= 1"-->
<!--                      [class.text-gray-700]="!username || currentStep() > 1"-->
<!--                      [class.bg-gray-400]="!username || currentStep() > 1"-->
<!--                      [class.cursor-not-allowed]="!username || currentStep() > 1"-->
<!--                      class="relative -ml-px inline-flex items-center gap-x-1.5 rounded-r-md px-3 py-2 text-sm font-medium focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-indigo-500 dark:focus:ring-indigo-400">-->
<!--                Next-->
<!--              </button>-->
<!--            </div>-->
<!--                  </span>-->
<!--                </div>-->
<!--              </li>-->

              <li class="relative">
                @if (currentStep() < 6) {
                  <div aria-hidden="true" class="absolute top-4 left-4 mt-0.5 -ml-px h-full w-0.5"
                       [class.bg-indigo-600]="currentStep() > 2"
                       [class.dark:bg-indigo-500]="currentStep() > 2"
                       [class.bg-gray-300]="currentStep() <= 2"
                       [class.dark:bg-gray-700]="currentStep() <= 2"></div>
                }
                <!-- Step 2: Install Directory -->
                <div class="group relative flex items-start">
                  <span class="flex h-9 items-center">
                    <span class="relative z-10 flex size-8 items-center justify-center rounded-full"
                          [class.bg-indigo-600]="currentStep() > 2"
                          [class.group-hover:bg-indigo-800]="currentStep() > 2"
                          [class.dark:bg-indigo-500]="currentStep() > 2"
                          [class.dark:group-hover:bg-indigo-600]="currentStep() > 2"
                          [class.border-2]="currentStep() <= 2"
                          [class.border-indigo-600]="currentStep() === 2"
                          [class.bg-white]="currentStep() === 2"
                          [class.dark:border-indigo-500]="currentStep() === 2"
                          [class.dark:bg-gray-900]="currentStep() === 2"
                          [class.border-gray-300]="currentStep() < 2"
                          [class.group-hover:border-gray-400]="currentStep() < 2"
                          [class.dark:border-white&#47;15]="currentStep() < 2"
                          [class.dark:group-hover:border-white&#47;25]="currentStep() < 2">
                      @if (currentStep() > 2) {
                        <svg viewBox="0 0 20 20" fill="currentColor" data-slot="icon"
                             aria-hidden="true" class="size-5 text-white">
                          <path
                            d="M16.704 4.153a.75.75 0 0 1 .143 1.052l-8 10.5a.75.75 0 0 1-1.127.075l-4.5-4.5a.75.75 0 0 1 1.06-1.06l3.894 3.893 7.48-9.817a.75.75 0 0 1 1.05-.143Z"
                            clip-rule="evenodd" fill-rule="evenodd"/>
                        </svg>
                      }
                      @if (currentStep() === 2) {
                        <span class="size-2.5 rounded-full bg-indigo-600 dark:bg-indigo-500"></span>
                      }
                      @if (currentStep() < 2) {
                        <span class="size-2.5 rounded-full bg-transparent group-hover:bg-gray-300 dark:group-hover:bg-white&#47;15"></span>
                      }
                    </span>
                  </span>
                  <span class="ml-4 flex min-w-0 flex-col">
                    <div [class.text-indigo-600]="currentStep() === 2"
                         [class.dark:text-indigo-400]="currentStep() === 2"
                         [class.text-gray-900]="currentStep() > 2"
                         [class.dark:text-white]="currentStep() > 2"
                         [class.text-gray-500]="currentStep() < 2"
                         [class.dark:text-gray-400]="currentStep() < 2"
                         class="mt-4 flex rounded-md shadow-sm">
              <div class="relative flex flex-grow items-stretch focus-within:z-10">
                <input type="text" name="installDirectory" id="installDirectory" [(ngModel)]="installDirectory"
                       class="block w-full rounded-none rounded-l-md border-0 py-1.5 text-gray-900 ring-1 ring-inset ring-gray-300 placeholder:text-gray-400 focus:ring-2 focus:ring-inset focus:ring-indigo-600 sm:text-sm sm:leading-6 dark:bg-gray-700 dark:text-white dark:ring-gray-600 dark:placeholder:text-gray-400 dark:focus:ring-indigo-500"
                       placeholder=" e.g., ~/lethean" [disabled]="currentStep() > 2">
                <button type="button" (click)="selectInstallDirectory()" [disabled]="currentStep() > 2"
                        class="relative -ml-px inline-flex items-center gap-x-1.5 rounded-r-md px-3 py-2 text-sm font-medium text-white bg-gray-600 hover:bg-gray-700 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-gray-500 dark:bg-gray-500 dark:hover:bg-gray-600 dark:focus:ring-gray-400">
                  Browse
                </button>
              </div>
              <button type="button" (click)="finishSetup()" [disabled]="!installDirectory || currentStep() > 2"
                      [class.text-white]="installDirectory && currentStep() <= 2"
                      [class.bg-indigo-600]="installDirectory && currentStep() <= 2"
                      [class.hover:bg-indigo-700]="installDirectory && currentStep() <= 2"
                      [class.dark:bg-indigo-500]="installDirectory && currentStep() <= 2"
                      [class.dark:hover:bg-indigo-600]="installDirectory && currentStep() <= 2"
                      [class.text-gray-700]="!installDirectory || currentStep() > 2"
                      [class.bg-gray-400]="!installDirectory || currentStep() > 2"
                      [class.cursor-not-allowed]="!installDirectory || currentStep() > 2"
                      class="relative -ml-px inline-flex items-center gap-x-1.5 rounded-r-md px-3 py-2 text-sm font-medium focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-indigo-500 dark:focus:ring-indigo-400">
                Next
              </button>
            </div>
                  </span>
                </div>
              </li>

            </ol>
          </nav>
        </div>

      </div>
    </div>
  `,
})
export class FullComponent {
  currentStep = signal(1);
  username = '';
  installDirectory = '~/lethean';

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
    if (this.currentStep() < 6) {
      this.currentStep.update(step => step + 1);
    }
  }

  async finishSetup(): Promise<void> {
    await Save()
    console.log('Setup finished!');
    console.log('Username:', this.username);
    console.log('Install Directory:', this.installDirectory);
    this.router.navigate(['/blockchain']);
  }
}
