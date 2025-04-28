import { Component } from '@angular/core';
import { Question, QuestionState } from '../../services/services';

@Component({
  selector: 'app-question',
  standalone: false,
  templateUrl: './question.component.html',
  styleUrl: './question.component.css',
})
export class QuestionComponent {
  question: Question = {
    title: 'Sample Question',
    statement: 'This is a sample question description.',
    limitations: {
      duration: 100,
      memory: 100,
    },
    state: QuestionState.QUESTION_STATE_PUBLISHED,
  };
}
