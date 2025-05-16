import { Component, Input, OnInit } from '@angular/core';
import { ManagerService } from '../../../services/manager.service';
import { ID, Question } from '../../../services/proto/services_pb';

@Component({
  selector: 'app-question-detail',
  standalone: false,
  templateUrl: './question-detail.component.html',
  styleUrl: './question-detail.component.css'
})
export class QuestionDetailComponent implements OnInit {
  @Input({ required: true })
  questionId!: string;

  question?: Question.AsObject;

  constructor(private readonly manager: ManagerService) {}

  ngOnInit() {
    this.manager
      .getQuestion(
        this.manager.create(new ID(), {
          value: this.questionId
        })
      )
      .then(res => {
        this.question = res.getQuestion()?.toObject();
      })
      .catch(err => {});
  }
}
