import { Component, ViewChild } from '@angular/core';
import { NewQuestionComponent } from './new-question/new-question.component';
import {
  ChangeQuestionStateRequest,
  Question
} from '../../services/proto/services_pb';
import { TabManagerComponent } from '../tab-manager/tab-manager.component';
import { ErrorHandlerService } from '../../services/error-handler.service';
import { ManagerService } from '../../services/manager.service';

@Component({
  selector: 'app-questions',
  standalone: false,
  templateUrl: './question-page.component.html',
  styleUrl: './question-page.component.css'
})
export class QuestionPageComponent {
  @ViewChild('editor')
  editor!: NewQuestionComponent;

  @ViewChild('tab')
  tab!: TabManagerComponent;

  constructor(
    private readonly errHandler: ErrorHandlerService,
    private readonly manager: ManagerService
  ) {}

  setQuestion(question?: Question.AsObject) {
    this.editor.setQuestion(question);
  }

  setTab(tab: number) {
    this.tab.activeTab = tab;
  }

  changeState(question: Question.AsObject) {
    this.manager
      .changeQuestionState(
        this.manager.create(new ChangeQuestionStateRequest(), {
          questionId: question.id,
          state: 2 - question.state
        })
      )
      .catch(err => {
        this.errHandler.handleError(err);
      });
  }
}
