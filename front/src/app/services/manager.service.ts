import { Injectable } from '@angular/core';
import { ManagerClient } from './proto/ServicesServiceClientPb';
import { CookieService } from './cookie.service';
import { Metadata } from 'grpc-web';
import { Router } from '@angular/router';

@Injectable({
  providedIn: 'root'
})
export class ManagerService extends ManagerClient {
  constructor(
    private readonly router: Router,
    private readonly cookie: CookieService
  ) {
    super('');
  }

  create<T>(t: T, properties: Record<string, any>): T {
    for (const key in properties) {
      (t as any)[`set${key.replace(/^[a-z]/, match => match.toUpperCase())}`](
        properties[key]
      );
    }
    return t;
  }

  getToken(): Metadata {
    return {
      authorization: `Bearer ${this.cookie.getCookie('token')}`
    };
  }

  reload() {
    this.router.navigateByUrl('/', { skipLocationChange: true }).then(() => {
      this.router.navigate([this.router.url]);
    });
  }
}
