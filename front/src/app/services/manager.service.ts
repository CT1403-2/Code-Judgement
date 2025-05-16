import { Injectable } from '@angular/core';
import { ManagerClient } from './proto/ServicesServiceClientPb';

@Injectable({
  providedIn: 'root'
})
export class ManagerService extends ManagerClient {
  constructor() {
    super('/api');
  }

  create<T>(t: T, properties: Record<string, any>): T {
    for (const key in properties) {
      (t as any)[`set${key.replace(/^[a-z]/, match => match.toUpperCase())}`](
        properties[key]
      );
    }
    return t;
  }
}
