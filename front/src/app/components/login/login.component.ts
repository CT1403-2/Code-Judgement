import { Component } from '@angular/core';
import { ManagerService } from '../../services/manager.service';
import { CookieService } from '../../services/cookie.service';
import {
  AuthenticationRequest,
  AuthenticationResponse
} from '../../services/proto/services_pb';
import { Router } from '@angular/router';
import { ErrorHandlerService } from '../../services/error-handler.service';

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
    private readonly router: Router,
    private readonly cookie: CookieService,
    private readonly errHandler: ErrorHandlerService,
    private readonly manager: ManagerService
  ) {}

  login(): void {
    this.handleAuth(
      this.manager.login(
        this.manager.create(new AuthenticationRequest(), {
          username: this.username,
          password: this.password
        })
      )
    );
  }

  signup(): void {
    this.handleAuth(
      this.manager.register(
        this.manager.create(new AuthenticationRequest(), {
          username: this.username,
          password: this.password
        })
      )
    );
  }

  handleAuth(authRes: Promise<AuthenticationResponse>) {
    authRes
      .then(res => {
        let date = new Date();
        date.setTime(date.getTime() + 24 * 60 * 60 * 1000);
        this.cookie.setCookie('token', res.getValue(), date);
        this.cookie.setCookie('role', `${res.getRole()}`, date);
        this.router.navigate(['questions']);
      })
      .catch(err => {
        this.errHandler.handleError(err);
      });
  }
}
