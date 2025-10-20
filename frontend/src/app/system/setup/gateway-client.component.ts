import { Component } from '@angular/core';
import { CommonModule } from '@angular/common';

@Component({
  selector: 'app-setup-gateway-client',
  standalone: true,
  imports: [CommonModule],
  template: `
    <div class="min-h-screen flex items-center justify-center bg-gray-50 py-12 px-4 sm:px-6 lg:px-8">
      <div class="max-w-md w-full space-y-8">
        <div>
          <h2 class="mt-6 text-center text-3xl font-extrabold text-gray-900">
            Gateway Client Setup
          </h2>
          <p class="mt-2 text-center text-sm text-gray-600">
            Configure your gateway client settings.
          </p>
        </div>
        <div class="mt-8 space-y-4">
          <button type="button" class="group relative w-full flex justify-center py-2 px-4 border border-transparent text-sm font-medium rounded-md text-white bg-green-600 hover:bg-green-700 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-green-500">
            Save Gateway Client Configuration
          </button>
        </div>
      </div>
    </div>
  `,
})
export class GatewayClientSetupComponent {}
