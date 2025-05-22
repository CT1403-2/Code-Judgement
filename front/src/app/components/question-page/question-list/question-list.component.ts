import { Component, EventEmitter, Input, Output } from '@angular/core';
import { Router } from '@angular/router';
import { ManagerService } from '../../../services/manager.service';
import {
  Filter,
  GetQuestionsRequest,
  ID,
  Question,
  Role
} from '../../../services/proto/services_pb';
import { ErrorHandlerService } from '../../../services/error-handler.service';

@Component({
  selector: 'app-question-list',
  standalone: false,
  templateUrl: './question-list.component.html',
  styleUrl: './question-list.component.css'
})
export class QuestionListComponent {
  questions!: Question.AsObject[];
  totalPageCount!: number;
  canChange: boolean = false;

  @Input()
  filterOwner: boolean = false;

  @Output()
  action = new EventEmitter<Question.AsObject>();

  constructor(
    private readonly router: Router,
    private readonly errHandler: ErrorHandlerService,
    private readonly manager: ManagerService
  ) {}

  gotoQuestion(question?: string) {
    this.router.navigate(['questions', question]);
  }

  fetchPage(page: number) {
    this.manager
      .getQuestions(
        this.manager.create(new GetQuestionsRequest(), {
          filtersList: [
            this.manager.create(new Filter(), {
              field: 'page',
              value: `${page}`
            }),
            this.manager.create(new Filter(), {
              field: 'owner',
              value: `${this.filterOwner}`
            })
          ]
        }),
        this.manager.getToken()
      )
      .then(res => {
        this.questions = res.toObject().questionsList;
        this.totalPageCount = res.getTotalPageSize();
      })
      .catch(err => {
        this.errHandler.handleError(err);
      });
    this.manager
      .getProfile(
        this.manager.create(new ID(), {
          value: ''
        }),
        this.manager.getToken()
      )
      .then(res => {
        this.canChange = res.getRole() > Role.ROLE_MEMBER;
      })
      .catch(err => {
        this.errHandler.handleError(err);
      });
  }
}
