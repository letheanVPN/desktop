import { Component } from '@angular/core';
import { CommonModule } from '@angular/common';

@Component({
  selector: 'app-setup-blockchain',
  standalone: true,
  imports: [CommonModule],
  template: `
    <div class="min-h-screen flex items-center justify-center bg-gray-50 py-12 px-4 sm:px-6 lg:px-8">
      <div class="max-w-md w-full space-y-8">
        <div>
          <h2 class="mt-6 text-center text-3xl font-extrabold text-gray-900">
            Blockchain Setup
          </h2>
          <p class="mt-2 text-center text-sm text-gray-600">
            Configure your blockchain settings.
          </p>
        </div>
        <div class="mt-8 space-y-4">
          <button type="button" class="group relative w-full flex justify-center py-2 px-4 border border-transparent text-sm font-medium rounded-md text-white bg-blue-600 hover:bg-blue-700 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-blue-500">
            Save Blockchain Configuration
          </button>
        </div>
      </div>
    </div>
  `,
})
export class BlockchainSetupComponent {}
