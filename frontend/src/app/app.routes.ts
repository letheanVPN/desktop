import { Routes } from '@angular/router';
import { ApplicationFrame } from '../frame/application.frame';
import { BlockchainComponent } from './blockchain/blockchain.component';
import { SystemTrayFrame } from '../frame/system-tray.frame';
import { DeveloperEditorComponent } from './developer/editor.component';
import { SetupComponent } from './system/setup.component';
import { FullComponent } from './system/setup/full.component';
import { BlockchainSetupComponent } from './system/setup/blockchain.component';
import { GatewayClientSetupComponent } from './system/setup/gateway-client.component';
import { SeedNodeSetupComponent } from './system/setup/seed-node.component';
import { MiningComponent } from './mining/mining.component';

export const routes: Routes = [
  { path: 'system-tray', component: SystemTrayFrame },
  { path: 'editor/monaco', component: DeveloperEditorComponent },
  {
    path: 'setup',
    component: SetupComponent,
    children: [
      { path: 'full', component: FullComponent },
      { path: 'blockchain', component: BlockchainSetupComponent },
      { path: 'gateway-client', component: GatewayClientSetupComponent },
      { path: 'seed-node', component: SeedNodeSetupComponent }
    ]
  },
  {
    path: '',
    component: ApplicationFrame,
    children: [
      { path: 'blockchain', component: BlockchainComponent },
      { path: 'dev/edit', component: DeveloperEditorComponent },
      { path: 'mining', component: MiningComponent },
      // Redirect empty path to a default view within the frame
      { path: '', redirectTo: 'blockchain', pathMatch: 'full' }
    ]
  }
];
