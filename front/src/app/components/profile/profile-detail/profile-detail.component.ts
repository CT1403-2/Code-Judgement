import { Component, Input, OnInit } from '@angular/core';
import { GetStatsResponse, Role } from '../../../services/services';
import * as chart from 'chart.js';
import { ManagerService } from '../../../services/manager.service';
import { grpc } from '@improbable-eng/grpc-web';

@Component({
  selector: 'app-profile-detail',
  standalone: false,
  templateUrl: './profile-detail.component.html',
  styleUrl: './profile-detail.component.css',
})
export class ProfileDetailComponent implements OnInit {
  roleTitles: { [key in Role]?: string } = {
    [Role.ROLE_UNKNOWN]: 'Unknown',
    [Role.ROLE_MEMBER]: 'Member',
    [Role.ROLE_ADMIN]: 'Admin',
    [Role.ROLE_SUPERUSER]: 'SuperUser',
  };

  @Input({ required: true })
  username!: string;

  stats!: GetStatsResponse;
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
            this.stats.solvedQuestions,
          ],
        },
      ],
    },
    options: {
      responsive: true,
      plugins: {
        legend: {
          position: 'top',
        },
      },
    },
  };

  constructor(private readonly manager: ManagerService) {}

  ngOnInit() {
    this.manager
      .GetProfile({
        value: this.username,
      })
      .then((res) => {
        this.role = res.role;
      })
      .catch((err) => {});
    this.manager
      .GetProfile({
        value: '',
      })
      .then((res) => {
        this.canChange = res.role > this.role;
      })
      .catch((err) => {});
    this.manager
      .GetStatsRequest({
        value: this.username,
      })
      .then((res) => {
        this.stats = res;
      })
      .catch((err) => {});
  }

  changeRole() {
    this.manager.ChangeRole({
      username: this.username,
      role: 3 - this.role,
    });
  }
}
