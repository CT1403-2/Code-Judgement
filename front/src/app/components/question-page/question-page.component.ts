import { Component } from '@angular/core';

@Component({
  selector: 'app-questions',
  standalone: false,
  templateUrl: './question-page.component.html',
  styleUrl: './question-page.component.css'
})
export class QuestionPageComponent {
  activeTab: number = 0
}
