import { Routes } from '@angular/router';
import { ApplicationFrame } from '../frame/application.frame';
import { BlockchainComponent } from './blockchain/blockchain.component';
import { SystemTrayFrame } from '../frame/system-tray.frame';
import { DeveloperEditorComponent } from './developer/editor.component';

export const routes: Routes = [
  { path: 'system-tray', component: SystemTrayFrame },
  { path: 'editor/monaco', component: DeveloperEditorComponent },
  {
    path: '',
    component: ApplicationFrame,
    children: [
      { path: 'blockchain', component: BlockchainComponent },
      { path: 'dev/edit', component: DeveloperEditorComponent },
      // Redirect empty path to a default view within the frame
      { path: '', redirectTo: 'blockchain', pathMatch: 'full' }
    ]
  }
];
