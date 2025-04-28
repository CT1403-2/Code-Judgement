import { Component } from '@angular/core';
import { QuestionState } from '../../../services/services';
import { ManagerService } from '../../../services/manager.service';

@Component({
  selector: 'app-new-question',
  standalone: false,
  templateUrl: './new-question.component.html',
  styleUrl: './new-question.component.css',
})
export class NewQuestionComponent {
  question = {
    title: '',
    statement: '',
    limitations: {
      duration: 0,
      memory: 0,
    },
    input: '',
    output: '',
    state: QuestionState.QUESTION_STATE_DRAFT,
  };

  constructor(private readonly manager: ManagerService) {}

  onSave(): void {
    this.manager.CreateQuestion(this.question).catch((err) => {});
  }
}
