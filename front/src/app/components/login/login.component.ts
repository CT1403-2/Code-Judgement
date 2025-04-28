import { Component } from '@angular/core';
import {ManagerService} from '../../services/manager.service';
import {CookieService} from '../../services/cookie.service';

@Component({
  selector: 'app-login',
  standalone: false,
  templateUrl: './login.component.html',
  styleUrl: './login.component.css',
})
export class LoginComponent {
  username: string = '';
  password: string = '';

  constructor(
    private readonly manager: ManagerService,
    private readonly cookie: CookieService
    ) {
  }

  login(): void {
    this.manager.Login({
      username: this.username,
      password: this.password
    }).then((res) => {
      let date = new Date();
      date.setTime(date.getTime() + (24 * 60 * 60 * 1000));
      this.cookie.setCookie('token', res.value, date);
    }).catch((err) => {

    });
  }

  signup(): void {
  }
}
