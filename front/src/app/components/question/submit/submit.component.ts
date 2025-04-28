import { Component } from '@angular/core';

@Component({
  selector: 'app-submit',
  standalone: false,
  templateUrl: './submit.component.html',
  styleUrl: './submit.component.css',
})
export class SubmitComponent {
  codeInput: string = '';
  fileContent: string = '';

  onFileSelected(event: Event): void {
    const input = event.target as HTMLInputElement;
    if (input.files && input.files.length > 0) {
      const file = input.files[0];
      const reader = new FileReader();
      reader.onload = (e) => {
        this.fileContent = e.target?.result as string;
        this.codeInput = this.fileContent; // Populate the textarea with file content
      };
      reader.readAsText(file);
    }
  }

  onSubmit(): void {
    const submission = this.codeInput.trim();
    if (submission) {
      console.log('Submitted code:', submission);
    } else {
      console.error('No code to submit.');
    }
  }
}
