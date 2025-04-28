import { Component, Input, OnInit } from '@angular/core';
import { GetStatsResponse } from '../../../services/services';
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
  @Input({ required: true })
  username!: string;

  stats!: GetStatsResponse;

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
      .GetStatsRequest({
        value: this.username,
      })
      .then((res) => {
        this.stats = res;
      })
      .catch((err) => {});
  }
}
