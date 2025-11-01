import {Component, CUSTOM_ELEMENTS_SCHEMA} from '@angular/core';
import { RouterOutlet } from '@angular/router';
import { setBasePath } from '@awesome.me/webawesome';
setBasePath('@awesome.me/webawesome/dist');
@Component({
  selector: 'app-root',
  standalone: true,
  imports: [RouterOutlet],
  schemas: [CUSTOM_ELEMENTS_SCHEMA],
  template: '<router-outlet></router-outlet>',
})
export class AppComponent {

}
