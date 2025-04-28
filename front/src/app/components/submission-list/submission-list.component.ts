import {Component, Input, OnInit} from '@angular/core';
import {Filter, Submission, SubmissionState} from '../../services/services';
import { Router } from '@angular/router';
import {ManagerService} from '../../services/manager.service';

@Component({
  selector: 'app-submission-list',
  standalone: false,
  templateUrl: './submission-list.component.html',
  styleUrl: './submission-list.component.css',
})
export class SubmissionListComponent implements OnInit {
  stateTitles: {[key in SubmissionState]?: string} = {
    [SubmissionState.SUBMISSION_STATE_UNKNOWN]: 'Unknown',
    [SubmissionState.SUBMISSION_STATE_PENDING]: 'Pending',
    [SubmissionState.SUBMISSION_STATE_JUDGING]: 'Judging',
    [SubmissionState.SUBMISSION_STATE_OK]: 'Ok',
    [SubmissionState.SUBMISSION_STATE_COMPILE_ERROR]: 'Compile Error',
    [SubmissionState.SUBMISSION_STATE_WRONG_ANSWER]: 'Wrong Answer',
    [SubmissionState.SUBMISSION_STATE_MEMORY_LIMIT_EXCEEDED]: 'Memory Limit Exceeded',
    [SubmissionState.SUBMISSION_STATE_TIME_LIMIT_EXCEEDED]: 'Time Limit Exceeded',
    [SubmissionState.SUBMISSION_STATE_RUNTIME_ERROR]: 'Runtime Error',
  }

  submissions!: Submission[];

  @Input()
  question?: string;

  constructor(
    private readonly router: Router,
    private readonly manager: ManagerService,
    ) {}

  gotoQuestion(question?: string) {
    this.router.navigate(['questions', question]);
  }

  ngOnInit() {
    let filters: Filter[] = [];
    if (this.question) {
      filters.push({
        field: 'questionId',
        value: this.question,
      });
    }
    this.manager.GetSubmissions({
      filters: filters,
    }).then((res) => {
      this.submissions = res.submissions.map((value) => {
        let val = value as any;
        val.stateTitle = this.stateTitles[value.state ?? 0];
        return val as Submission;
      });
    }).catch((err) => {

    });
  }
}
