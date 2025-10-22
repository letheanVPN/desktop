import { Component, OnInit } from '@angular/core';
import { MonacoEditorModule } from 'ngx-monaco-editor-v2';
import {FormsModule} from '@angular/forms';
import {Get} from "@lthn/core/config/service"
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
export class DeveloperEditorComponent implements OnInit {
  editorOptions = { theme: 'vs-dark', language: 'typescript' };
  code: string = "// Start Coding"; // Initialize with a default value

  async ngOnInit(): Promise<void> {
    try {
      const cfg = await Get();
      this.code = cfg ? JSON.stringify(cfg, null, 2) : "// Start Coding";
    } catch (error) {
      console.error("Error fetching configuration:", error);
      this.code = "// Error loading configuration.";
    }
  }
}
