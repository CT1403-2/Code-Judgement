import { Component, OnInit } from '@angular/core';
import { ActivatedRoute } from '@angular/router';

@Component({
  selector: 'app-question',
  standalone: false,
  templateUrl: './question.component.html',
  styleUrl: './question.component.css'
})
export class QuestionComponent implements OnInit {
  question!: string;

  constructor(private readonly route: ActivatedRoute) {}

  ngOnInit() {
    this.route.params.subscribe(params => {
      this.question = params['id'];
    });
  }
}
