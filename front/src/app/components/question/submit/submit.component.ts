import { Component, Input } from '@angular/core';
import { ManagerService } from '../../../services/manager.service';
import { Submission, SubmitRequest } from '../../../services/proto/services_pb';
import { ErrorHandlerService } from '../../../services/error-handler.service';

@Component({
  selector: 'app-submit',
  standalone: false,
  templateUrl: './submit.component.html',
  styleUrl: './submit.component.css'
})
export class SubmitComponent {
  codeInput: string = '';
  file?: File;

  @Input({ required: true })
  question!: string;

  constructor(
    private readonly errHandler: ErrorHandlerService,
    private readonly manager: ManagerService
  ) {}

  onSubmit(): void {
    if (this.file) {
      const reader = new FileReader();
      reader.onload = () => {
        this.submit(new Uint8Array(reader.result as ArrayBuffer));
      };
      reader.readAsArrayBuffer(this.file);
    } else {
      this.submit(new TextEncoder().encode(this.codeInput));
    }
  }

  submit(codeData: Uint8Array) {
    this.manager
      .submit(
        this.manager.create(new SubmitRequest(), {
          submission: this.manager.create(new Submission(), {
            questionId: this.question,
            code: codeData
          })
        })
      )
      .then(() => {
        this.manager.reload();
      })
      .catch(err => {
        this.errHandler.handleError(err);
      });
  }
}
