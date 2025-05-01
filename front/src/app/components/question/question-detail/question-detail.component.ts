import { Component, Input, OnInit } from '@angular/core';
import { Question } from '../../../services/services';
import { ManagerService } from '../../../services/manager.service';

@Component({
  selector: 'app-question-detail',
  standalone: false,
  templateUrl: './question-detail.component.html',
  styleUrl: './question-detail.component.css',
})
export class QuestionDetailComponent implements OnInit {
  @Input({ required: true })
  questionId!: string;

  question?: Question;

  constructor(private readonly manager: ManagerService) {}

  ngOnInit() {
    this.manager
      .GetQuestion({
        value: this.questionId,
      })
      .then((res) => {
        this.question = res.question;
      })
      .catch((err) => {});
  }
}
