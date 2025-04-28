import {Component, OnInit} from '@angular/core';
import {ActivatedRoute, Router} from '@angular/router';

@Component({
  selector: 'app-error',
  standalone: false,
  templateUrl: './error.component.html',
  styleUrl: './error.component.css',
})
export class ErrorComponent implements OnInit {
  code!: number;

  messages: { [key: number]: string } = {
    404: 'Page Not Found',
    403: 'Access Denied',
  };

  constructor(
    private readonly route: ActivatedRoute,
    private readonly router: Router
    ) {}

  ngOnInit(): void {
    this.route.params.subscribe((params) => {
      this.code = params['id'];
      if (this.messages[this.code] == undefined) {
        this.router.navigate(['/error', 404]);
      }
    });
  }
}
