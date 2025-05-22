import { Component } from '@angular/core';
import {
  ErrorHandlerService,
  ErrorMessage
} from '../../services/error-handler.service';

@Component({
  selector: 'app-error-dialog',
  standalone: false,
  templateUrl: './error-dialog.component.html',
  styleUrl: './error-dialog.component.css'
})
export class ErrorDialogComponent {
  messages: ErrorMessage[] = [];

  get isEmpty(): boolean {
    return this.messages.length === 0;
  }

  constructor(private errHandler: ErrorHandlerService) {
    this.errHandler.error$.subscribe(msg => this.messages.push(msg));
  }

  close() {
    while (!this.isEmpty) {
      this.messages.pop()?.resolve();
    }
  }
}
