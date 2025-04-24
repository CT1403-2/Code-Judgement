import { Component } from '@angular/core';

@Component({
  selector: 'app-login',
  standalone: false,
  templateUrl: './login.component.html',
  styleUrl: './login.component.css'
})
export class LoginComponent {
  email: string = '';
  password: string = '';

  onSubmit(): void {
    console.log('Email:', this.email);
    console.log('Password:', this.password);
    // Add authentication logic here
  }
}
