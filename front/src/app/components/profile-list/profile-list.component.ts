import { Component, OnInit } from '@angular/core';
import { Router } from '@angular/router';
import { ManagerService } from '../../services/manager.service';
import { GetProfilesRequest, ID } from '../../services/proto/services_pb';
import * as grpcWeb from 'grpc-web';

@Component({
  selector: 'app-profile-list',
  standalone: false,
  templateUrl: './profile-list.component.html',
  styleUrl: './profile-list.component.css'
})
export class ProfileListComponent implements OnInit {
  currentUser!: string;
  profiles!: { username: string }[];

  constructor(
    private readonly router: Router,
    private readonly manager: ManagerService
  ) {}

  gotoMyProfile(tab: number) {
    if (tab == 1) {
      this.gotoProfile();
    }
  }

  gotoProfile(user: string = this.currentUser) {
    this.router.navigate(['profiles', user]);
  }

  ngOnInit() {
    this.manager
      .getProfile(
        this.manager.create(new ID(), {
          value: ''
        })
      )
      .then(res => {
        this.currentUser = res.getUsername();
      })
      .catch((err: grpcWeb.RpcError) => {
        if (err.code == grpcWeb.StatusCode.UNAUTHENTICATED) {
          this.router.navigate(['']);
        }
      });
    this.manager
      .getProfiles(
        this.manager.create(new GetProfilesRequest(), {
          filtersList: []
        })
      )
      .then(res => {
        this.profiles = res.getUsernamesList().map(user => {
          return { username: user };
        });
      })
      .catch(err => {
      });
  }
}
