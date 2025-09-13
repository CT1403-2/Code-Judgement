import { Component, OnInit } from '@angular/core';
import { ActivatedRoute, Router } from '@angular/router';

@Component({
  selector: 'app-error-page',
  standalone: false,
  templateUrl: './error-page.component.html',
  styleUrl: './error-page.component.css'
})
export class ErrorPageComponent implements OnInit {
  code!: number;

  messages: { [key: number]: string } = {
    404: 'Page Not Found',
    403: 'Access Denied'
  };

  constructor(
    private readonly route: ActivatedRoute,
    private readonly router: Router
  ) {}

  ngOnInit(): void {
    this.route.params.subscribe(params => {
      this.code = params['id'];
      if (this.messages[this.code] == undefined) {
        this.router.navigate(['/error', 404]);
      }
    });
  }
}
