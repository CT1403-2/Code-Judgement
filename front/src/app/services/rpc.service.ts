import { Injectable } from '@angular/core';

interface Rpc {
  request(service: string, method: string, data: Uint8Array): Promise<Uint8Array>;
}

@Injectable({
  providedIn: 'root'
})
export class RpcService implements Rpc {
  private readonly baseUrl = "";

  constructor() {
  }

  async request(service: string, method: string, data: Uint8Array): Promise<Uint8Array> {
    const url = `${this.baseUrl}/${service}/${method}`;

    return new Promise((resolve, reject) => {
      fetch(url, {
        method: 'POST',
        body: data,
        headers: {
          'Content-Type': 'application/octet-stream',
        },
      }).then((response) => {
        if (!response.ok) {
          reject(new Error(`Request failed with status: ${response.status}`));
        }

        response.arrayBuffer().then((responseData) => {
          resolve(new Uint8Array(responseData));
        }).catch((error) => {
          reject(error);
        });
      }).catch((error) => {
        reject(error);
      });
    });
  }
}
