import {Component, CUSTOM_ELEMENTS_SCHEMA, OnInit} from '@angular/core';
import { CommonModule } from '@angular/common';

@Component({
  selector: 'app-mining',
  standalone: true,
  imports: [CommonModule],
  schemas: [CUSTOM_ELEMENTS_SCHEMA],
  template: `
      <!-- Sticky search header -->
      <div slot="subheader" class=" top-0 z-40 flex h-16 shrink-0 items-center gap-x-6 border-b border-gray-200 bg-white  px-4 shadow-xs sm:px-6 lg:px-8 dark:border-white/5 dark:bg-gray-900 dark:shadow-none">
        <nav class="flex overflow-x-auto border-b border-gray-200 py-4 dark:border-white/10">
          <div class="pt-2 pr-2">
            <div class="flex-none rounded-full bg-green-500/10 text-green-500 dark:bg-green-400/10 dark:text-green-400">
              <div class="size-2 rounded-full bg-current"></div>
            </div>
          </div>
          <ul role="list" class="flex min-w-full flex-none gap-x-6 px-1 text-sm/6 font-semibold text-gray-500 sm:px-1 lg:px-1 dark:text-gray-400">
            <li>
              <a href="#" class="text-indigo-600 dark:text-indigo-400">Overview</a>
            </li>
            <li>
              <a href="#" class="">Stats</a>
            </li>
            <li>
              <a href="#" class="">Manage</a>
            </li>
            <li>
              <a href="#" class="">Settings</a>
            </li>
            <li>
              <a href="#" class="">Notifications</a>
            </li>
          </ul>
        </nav>
      </div>
      <main class="flex-1">
        <div class="grid grid-cols-1 gap-4 sm:grid-cols-2 lg:grid-cols-3 p-4">
          <!-- Current mining stats -->
          <div class="bg-white dark:bg-gray-800 shadow-md rounded-lg p-4">
            <h2 class="text-lg font-semibold text-gray-900 dark:text-white mb-4 flex items-center">
              <i class="fa-regular fa-chart-line mr-2"></i> Current Mining Stats
            </h2>
            <div class="grid grid-cols-2 gap-4">
              <div>
                <p class="text-sm text-gray-500 dark:text-gray-400">Hashrate</p>
                <p class="text-2xl font-bold text-gray-900 dark:text-white">12.5 KH/s</p>
              </div>
              <div>
                <p class="text-sm text-gray-500 dark:text-gray-400">Temperature</p>
                <p class="text-2xl font-bold text-gray-900 dark:text-white">65°C</p>
              </div>
              <div>
                <p class="text-sm text-gray-500 dark:text-gray-400">Uptime</p>
                <p class="text-2xl font-bold text-gray-900 dark:text-white">2h 15m</p>
              </div>
              <div>
                <p class="text-sm text-gray-500 dark:text-gray-400">Blocks Found</p>
                <p class="text-2xl font-bold text-gray-900 dark:text-white">3</p>
              </div>
            </div>
          </div>

          <!-- Mining config quick run -->
          <div class="bg-white dark:bg-gray-800 shadow-md rounded-lg p-4">
            <h2 class="text-lg font-semibold text-gray-900 dark:text-white mb-4 flex items-center">
              <i class="fa-regular fa-rocket-launch mr-2"></i> Mining Config Quick Run
            </h2>
            <div class="space-y-4">
              <div>
                <label for="mining-profile" class="block text-sm font-medium text-gray-700 dark:text-gray-300">Select Profile</label>
                <select id="mining-profile" name="mining-profile" class="mt-1 block w-full pl-3 pr-10 py-2 text-base border-gray-300 focus:outline-none focus:ring-indigo-500 focus:border-indigo-500 sm:text-sm rounded-md dark:bg-gray-700 dark:border-gray-600 dark:placeholder-gray-400 dark:text-white">
                  <option>Default</option>
                  <option>High Performance</option>
                  <option>Low Power</option>
                </select>
              </div>
              <button type="button" class="w-full inline-flex justify-center items-center px-4 py-2 border border-transparent text-sm font-medium rounded-md shadow-sm text-white bg-indigo-600 hover:bg-indigo-700 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-indigo-500">
                <i class="fad fa-play mr-2"></i> Start Mining
              </button>
            </div>
            <div class="mt-4 text-center">
              <a href="#" class="text-indigo-600 dark:text-indigo-400 hover:text-indigo-500 text-sm">
                <i class="fa-regular fa-plus-circle mr-1"></i> Add new run profile
              </a>
            </div>
          </div>

          <!-- Installed software version -->
          <div class="bg-white dark:bg-gray-800 shadow-md rounded-lg p-4">
            <h2 class="text-lg font-semibold text-gray-900 dark:text-white mb-4 flex items-center">
              <i class="fa-regular fa-info-circle mr-2"></i> Installed Software Version
            </h2>
            <div class="text-center">
              <p class="text-sm text-gray-500 dark:text-gray-400">XMRig</p>
              <p class="text-3xl font-bold text-gray-900 dark:text-white">v6.18.0</p>
              <a href="#" class="mt-4 inline-block text-indigo-600 dark:text-indigo-400 hover:text-indigo-500 text-sm">
                <i class="fa-regular fa-sync-alt mr-1"></i> Check for update
              </a>
            </div>
          </div>
        </div>
      </main>

  `,
  styles: []
})
export class MiningComponent  {

  constructor() { }


}
