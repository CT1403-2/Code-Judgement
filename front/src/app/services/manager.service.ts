import { Injectable } from '@angular/core';
import {RpcService} from './rpc.service';
import {ManagerClientImpl} from './services';

@Injectable({
  providedIn: 'root'
})
export class ManagerService extends ManagerClientImpl {
  constructor(rpc: RpcService) {
    super(rpc, {
      service: "",
    });
  }
}
