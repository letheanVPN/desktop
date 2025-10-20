import { Component, OnDestroy, OnInit } from '@angular/core';
import { CommonModule } from '@angular/common';
import { RouterOutlet } from '@angular/router';

@Component({
  selector: 'blank-frame',
  standalone: true,
  imports: [CommonModule, RouterOutlet],
  template: `
    <router-outlet></router-outlet>
  `,
})
export class BlankFrame implements OnInit, OnDestroy {
  ngOnInit(): void {
    // Initialization logic for blank frame
  }

  ngOnDestroy(): void {
    // Cleanup logic for blank frame
  }
}
