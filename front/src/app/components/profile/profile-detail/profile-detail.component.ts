import { Component, Input, OnInit } from '@angular/core';
import * as chart from 'chart.js';
import { ManagerService } from '../../../services/manager.service';
import {
  ChangeRoleRequest,
  GetStatsResponse,
  ID,
  Role
} from '../../../services/proto/services_pb';
import { ErrorHandlerService } from '../../../services/error-handler.service';

@Component({
  selector: 'app-profile-detail',
  standalone: false,
  templateUrl: './profile-detail.component.html',
  styleUrl: './profile-detail.component.css'
})
export class ProfileDetailComponent implements OnInit {
  roleTitles: { [key in Role]?: string } = {
    [Role.ROLE_UNKNOWN]: 'Unknown',
    [Role.ROLE_MEMBER]: 'Member',
    [Role.ROLE_ADMIN]: 'Admin',
    [Role.ROLE_SUPERUSER]: 'SuperUser'
  };

  @Input({ required: true })
  username!: string;

  stats: GetStatsResponse.AsObject = {
    solvedQuestions: 0,
    triedQuestions: 0
  };
  role!: Role;
  canChange: boolean = false;

  config: chart.ChartConfiguration = {
    type: 'doughnut',
    data: {
      labels: ['Unsolved Questions', 'Solved Questions'],
      datasets: [
        {
          data: [
            this.stats.triedQuestions - this.stats.solvedQuestions,
            this.stats.solvedQuestions
          ]
        }
      ]
    },
    options: {
      responsive: true,
      plugins: {
        legend: {
          position: 'top'
        }
      }
    }
  };

  constructor(
    private readonly errHandler: ErrorHandlerService,
    private readonly manager: ManagerService
  ) {}

  ngOnInit() {
    this.manager
      .getProfile(
        this.manager.create(new ID(), {
          value: this.username
        }),
        this.manager.getToken()
      )
      .then(res => {
        this.role = res.getRole();
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
        this.canChange = res.getRole() > this.role;
      })
      .catch(err => {
        this.errHandler.handleError(err);
      });
    this.manager
      .getStatsRequest(
        this.manager.create(new ID(), {
          value: this.username
        }),
        this.manager.getToken()
      )
      .then(res => {
        this.stats = res.toObject();
      })
      .catch(err => {
        this.errHandler.handleError(err);
      });
  }

  changeRole() {
    this.manager
      .changeRole(
        this.manager.create(new ChangeRoleRequest(), {
          username: this.username,
          role: 3 - this.role
        })
      )
      .then(() => {
        this.manager.reload();
      })
      .catch(err => {
        this.errHandler.handleError(err);
      });
  }
}
