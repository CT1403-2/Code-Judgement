import { Injectable } from '@angular/core';
import { Router } from '@angular/router';
import * as grpcWeb from 'grpc-web';
import { Subject } from 'rxjs';

export interface ErrorMessage {
  message: string;
  resolve: () => void;
}

@Injectable({
  providedIn: 'root'
})
export class ErrorHandlerService {
  private errorSubject = new Subject<ErrorMessage>();
  error$ = this.errorSubject.asObservable();
  constructor(private readonly router: Router) {}

  showError(message: string): Promise<void> {
    return new Promise(resolve => {
      this.errorSubject.next({
        message,
        resolve
      });
    });
  }

  handleError(err: any) {
    if (err instanceof grpcWeb.RpcError) {
      this.showError(err.message).then(() => {
        switch (err.code) {
          case grpcWeb.StatusCode.UNAUTHENTICATED:
            this.router.navigate(['']);
            break;
          case grpcWeb.StatusCode.PERMISSION_DENIED:
            this.router.navigate(['error', '403']);
            break;;
          case grpcWeb.StatusCode.NOT_FOUND:
            this.router.navigate(['error', '40r']);
            break;
        }
      });
    }
  }
}
