import { Component, Input, OnInit } from '@angular/core';
import * as chart from 'chart.js';
import { ManagerService } from '../../../services/manager.service';
import {
  ChangeRoleRequest,
  GetStatsResponse,
  ID,
  Role
} from '../../../services/proto/services_pb';

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

  stats!: GetStatsResponse.AsObject;
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

  constructor(private readonly manager: ManagerService) {}

  ngOnInit() {
    this.manager
      .getProfile(
        this.manager.create(new ID(), {
          value: this.username
        })
      )
      .then(res => {
        this.role = res.getRole();
      })
      .catch(err => {});
    this.manager
      .getProfile(
        this.manager.create(new ID(), {
          value: ''
        })
      )
      .then(res => {
        this.canChange = res.getRole() > this.role;
      })
      .catch(err => {});
    this.manager
      .getStatsRequest(
        this.manager.create(new ID(), {
          value: this.username
        })
      )
      .then(res => {
        this.stats = res.toObject();
      })
      .catch(err => {});
  }

  changeRole() {
    this.manager
      .changeRole(
        this.manager.create(new ChangeRoleRequest(), {
          username: this.username,
          role: 3 - this.role
        })
      )
      .catch(err => {});
  }
}
