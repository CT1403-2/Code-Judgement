import { Component, Input } from '@angular/core';
import { ManagerService } from '../../../services/manager.service';
import { Submission, SubmitRequest } from '../../../services/proto/services_pb';

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

  constructor(private readonly manager: ManagerService) {}

  onSubmit(): void {
    if (this.file) {
      const reader = new FileReader();
      reader.onload = () => {
        const fileData = new Uint8Array(reader.result as ArrayBuffer);
        this.manager
          .submit(
            this.manager.create(new SubmitRequest(), {
              submission: this.manager.create(new Submission(), {
                questionId: this.question,
                code: fileData
              })
            })
          )
          .catch(err => {});
      };
      reader.readAsArrayBuffer(this.file);
    } else {
      const codeData = new TextEncoder().encode(this.codeInput);
      this.manager
        .submit(
          this.manager.create(new SubmitRequest(), {
            submission: this.manager.create(new Submission(), {
              questionId: this.question,
              code: codeData
            })
          })
        )
        .catch(err => {});
    }
  }
}
