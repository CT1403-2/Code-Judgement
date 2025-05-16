import { Component, OnInit } from '@angular/core';
import { Router } from '@angular/router';
import { ManagerService } from '../../../services/manager.service';
import {GetQuestionsRequest, Question} from '../../../services/proto/services_pb';

@Component({
  selector: 'app-question-list',
  standalone: false,
  templateUrl: './question-list.component.html',
  styleUrl: './question-list.component.css'
})
export class QuestionListComponent implements OnInit {
  questions!: Question.AsObject[];

  constructor(
    private readonly router: Router,
    private readonly manager: ManagerService
  ) {}

  gotoQuestion(question?: string) {
    this.router.navigate(['questions', question]);
  }

  ngOnInit() {
    this.manager
      .getQuestions(this.manager.create(new GetQuestionsRequest(), {
        filtersList: []
      }))
      .then(res => {
        this.questions = res.toObject().questionsList;
      })
      .catch(err => {});
  }
}
