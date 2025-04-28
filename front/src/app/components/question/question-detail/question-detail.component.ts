import { Component, Input } from '@angular/core';
import { Question } from '../../../services/services';

@Component({
  selector: 'app-question-detail',
  standalone: false,
  templateUrl: './question-detail.component.html',
  styleUrl: './question-detail.component.css',
})
export class QuestionDetailComponent {
  @Input({ required: true })
  question!: Question;
}
