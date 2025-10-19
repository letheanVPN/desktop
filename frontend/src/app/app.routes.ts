import { Routes } from '@angular/router';
import { ApplicationFrame } from '../frame/application.frame';
import { BlockchainComponent } from './blockchain/blockchain.component';
import { SystemTrayFrame } from '../frame/system-tray.frame';

export const routes: Routes = [
  { path: 'system-tray', component: SystemTrayFrame },
  {
    path: '',
    component: ApplicationFrame,
    children: [
      { path: 'blockchain', component: BlockchainComponent },
      // Redirect empty path to a default view within the frame
      { path: '', redirectTo: 'blockchain', pathMatch: 'full' }
    ]
  }
];
