import { Component, Input } from '@angular/core';
import { Router } from '@angular/router';
import { ManagerService } from '../../services/manager.service';
import {
  Filter,
  GetSubmissionsRequest,
  Submission,
  SubmissionState
} from '../../services/proto/services_pb';
import { ErrorHandlerService } from '../../services/error-handler.service';

@Component({
  selector: 'app-submission-list',
  standalone: false,
  templateUrl: './submission-list.component.html',
  styleUrl: './submission-list.component.css'
})
export class SubmissionListComponent {
  stateTitles: { [key in SubmissionState]?: string } = {
    [SubmissionState.SUBMISSION_STATE_UNKNOWN]: 'Unknown',
    [SubmissionState.SUBMISSION_STATE_PENDING]: 'Pending',
    [SubmissionState.SUBMISSION_STATE_JUDGING]: 'Judging',
    [SubmissionState.SUBMISSION_STATE_OK]: 'Ok',
    [SubmissionState.SUBMISSION_STATE_COMPILE_ERROR]: 'Compile Error',
    [SubmissionState.SUBMISSION_STATE_WRONG_ANSWER]: 'Wrong Answer',
    [SubmissionState.SUBMISSION_STATE_MEMORY_LIMIT_EXCEEDED]:
      'Memory Limit Exceeded',
    [SubmissionState.SUBMISSION_STATE_TIME_LIMIT_EXCEEDED]:
      'Time Limit Exceeded',
    [SubmissionState.SUBMISSION_STATE_RUNTIME_ERROR]: 'Runtime Error'
  };

  submissions!: Submission.AsObject[];
  totalPageCount!: number;

  @Input({ required: true })
  filterType!: 'questionId' | 'username';

  @Input({ required: true })
  filterValue!: string;

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
      .getSubmissions(
        this.manager.create(new GetSubmissionsRequest(), {
          filtersList: [
            this.manager.create(new Filter(), {
              field: 'page',
              value: `${page}`
            }),
            this.manager.create(new Filter(), {
              field: this.filterType,
              value: this.filterValue
            })
          ]
        }),
        this.manager.getToken()
      )
      .then(res => {
        this.submissions = res.getSubmissionsList().map(value => {
          let val = value.toObject() as any;
          val.stateTitle = this.stateTitles[
            value.hasState() ? value.getState() : 0
          ];
          return val as Submission.AsObject;
        });
        this.totalPageCount = res.getTotalPageSize();
      })
      .catch(err => {
        this.errHandler.handleError(err);
      });
  }
}
