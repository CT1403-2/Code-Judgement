import { Component } from '@angular/core';
import { QuestionState } from '../../../services/services';

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

  onSave(): void {
    if (this.question.title && this.question.statement) {
      console.log('Question saved:', this.question);
    } else {
      console.error('Please fill in all required fields.');
    }
  }
}
