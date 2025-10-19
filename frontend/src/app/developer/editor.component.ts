import { Component } from '@angular/core';
import { MonacoEditorModule } from 'ngx-monaco-editor-v2';
import {FormsModule} from '@angular/forms';

@Component({
  selector: 'dev-edit',
  standalone: true,
  imports: [MonacoEditorModule, FormsModule],
  template: `
    <div style="height: 100vh;">
      <ngx-monaco-editor style="height: 100%;" [options]="editorOptions" [(ngModel)]="code"></ngx-monaco-editor>
    </div>
  `,
})
export class DeveloperEditorComponent {
  editorOptions = { theme: 'vs-dark', language: 'typescript' };
  code = '// Start coding...';
}
