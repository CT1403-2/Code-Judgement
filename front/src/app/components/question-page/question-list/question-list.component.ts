import { Component } from '@angular/core';
import { Question, QuestionState } from '../../../services/services';
import { Router } from '@angular/router';

@Component({
  selector: 'app-question-list',
  standalone: false,
  templateUrl: './question-list.component.html',
  styleUrl: './question-list.component.css',
})
export class QuestionListComponent {
  questions: Question[] = [
    {
      id: '1',
      title: 'Sample Question',
      statement: 'This is a sample question description.',
      limitations: {
        duration: 100,
        memory: 100,
      },
      state: QuestionState.QUESTION_STATE_PUBLISHED,
    },
  ];

  constructor(private readonly router: Router) {}

  gotoQuestion(question?: string) {
    this.router.navigate(['questions', question]);
  }
}
