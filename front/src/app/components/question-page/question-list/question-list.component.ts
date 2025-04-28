import { Component, OnInit } from '@angular/core';
import { Question, QuestionState } from '../../../services/services';
import { Router } from '@angular/router';
import { ManagerService } from '../../../services/manager.service';

@Component({
  selector: 'app-question-list',
  standalone: false,
  templateUrl: './question-list.component.html',
  styleUrl: './question-list.component.css',
})
export class QuestionListComponent implements OnInit {
  questions!: Question[];

  constructor(
    private readonly router: Router,
    private readonly manager: ManagerService,
  ) {}

  gotoQuestion(question?: string) {
    this.router.navigate(['questions', question]);
  }

  ngOnInit() {
    this.manager
      .GetQuestions({
        filters: [],
      })
      .then((res) => {
        this.questions = res.questions;
      })
      .catch((err) => {});
  }
}
