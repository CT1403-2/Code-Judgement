import { Component, Input } from '@angular/core';
import { GetStatsResponse } from '../../../services/services';
import * as chart from 'chart.js';

@Component({
  selector: 'app-profile-detail',
  standalone: false,
  templateUrl: './profile-detail.component.html',
  styleUrl: './profile-detail.component.css',
})
export class ProfileDetailComponent {
  @Input({ required: true })
  username!: string;

  stats: GetStatsResponse = {
    triedQuestions: 10,
    solvedQuestions: 5,
  };

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
}
