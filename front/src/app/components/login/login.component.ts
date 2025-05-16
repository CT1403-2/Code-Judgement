import { Component } from '@angular/core';
import { ManagerService } from '../../services/manager.service';
import { CookieService } from '../../services/cookie.service';
import { AuthenticationRequest } from '../../services/proto/services_pb';

@Component({
  selector: 'app-login',
  standalone: false,
  templateUrl: './login.component.html',
  styleUrl: './login.component.css'
})
export class LoginComponent {
  username: string = '';
  password: string = '';

  constructor(
    private readonly manager: ManagerService,
    private readonly cookie: CookieService
  ) {}

  login(): void {
    this.manager
      .login(
        this.manager.create(new AuthenticationRequest(), {
          username: this.username,
          password: this.password
        })
      )
      .then(res => {
        let date = new Date();
        date.setTime(date.getTime() + 24 * 60 * 60 * 1000);
        this.cookie.setCookie('token', res.getValue(), date);
        this.cookie.setCookie('role', `${res.getRole()}`, date);
      })
      .catch(err => {});
  }

  signup(): void {
    this.manager
      .register(
        this.manager.create(new AuthenticationRequest(), {
          username: this.username,
          password: this.password
        })
      )
      .then(res => {
        let date = new Date();
        date.setTime(date.getTime() + 24 * 60 * 60 * 1000);
        this.cookie.setCookie('token', res.getValue(), date);
        this.cookie.setCookie('role', `${res.getRole()}`, date);
      })
      .catch(err => {});
  }
}
