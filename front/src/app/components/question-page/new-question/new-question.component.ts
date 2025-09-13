import { Component } from '@angular/core';
import { ManagerService } from '../../../services/manager.service';
import {
  ID,
  Limitations,
  Question,
  QuestionState
} from '../../../services/proto/services_pb';
import { ErrorHandlerService } from '../../../services/error-handler.service';

@Component({
  selector: 'app-new-question',
  standalone: false,
  templateUrl: './new-question.component.html',
  styleUrl: './new-question.component.css'
})
export class NewQuestionComponent {
  emptyQuestion: Question.AsObject = {
    title: '',
    statement: '',
    limitations: {
      duration: 0,
      memory: 0
    },
    input: '',
    output: '',
    owner: '',
    state: QuestionState.QUESTION_STATE_DRAFT
  };

  question!: any;

  constructor(
    private readonly errHandler: ErrorHandlerService,
    private readonly manager: ManagerService
  ) {}

  onSave(): void {
    const question = this.manager.create(new Question(), {
      title: this.question.title,
      statement: this.question.statement,
      limitations: this.manager.create(new Limitations(), {
        duration: this.question.limitation.duration,
        memory: this.question.limitation.memory
      }),
      input: this.question.input,
      output: this.question.output,
      owner: this.question.owner,
      state: QuestionState.QUESTION_STATE_DRAFT
    });
    this.handleQuestion(
      question.hasId()
        ? this.manager.editQuestion(question)
        : this.manager.createQuestion(question)
    );
  }

  handleQuestion(questionRes: Promise<any>) {
    questionRes
      .then(() => {
        this.manager.reload();
      })
      .catch(err => {
        this.errHandler.handleError(err);
      });
  }

  setQuestion(question?: Question.AsObject) {
    this.manager
      .getProfile(
        this.manager.create(new ID(), {
          value: ''
        }),
        this.manager.getToken()
      )
      .then(res => {
        const owner = res.getUsername();
        if (question) {
          if (question.owner != owner) {
            this.errHandler.showError(
              "you can't edit other peoples questions."
            );
            this.manager.reload();
          } else {
            this.question = structuredClone(question);
          }
        } else {
          this.question = structuredClone(this.emptyQuestion);
          this.question.owner = owner;
        }
      })
      .catch(err => {
        this.errHandler.handleError(err);
      });
  }
}
