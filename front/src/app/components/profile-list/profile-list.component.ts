import { Component, OnInit } from '@angular/core';
import { Router } from '@angular/router';
import { ManagerService } from '../../services/manager.service';
import {
  Filter,
  GetProfilesRequest,
  ID
} from '../../services/proto/services_pb';
import { ErrorHandlerService } from '../../services/error-handler.service';

@Component({
  selector: 'app-profile-list',
  standalone: false,
  templateUrl: './profile-list.component.html',
  styleUrl: './profile-list.component.css'
})
export class ProfileListComponent implements OnInit {
  currentUser!: string;
  profiles!: { username: string }[];
  totalPageCount!: number;

  constructor(
    private readonly router: Router,
    private readonly errHandler: ErrorHandlerService,
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
        }),
        this.manager.getToken()
      )
      .then(res => {
        this.currentUser = res.getUsername();
      })
      .catch(err => {
        this.errHandler.handleError(err);
      });
  }

  fetchPage(page: number) {
    this.manager
      .getProfiles(
        this.manager.create(new GetProfilesRequest(), {
          filtersList: [
            this.manager.create(new Filter(), {
              field: 'page',
              value: `${page}`
            })
          ]
        })
      )
      .then(res => {
        this.profiles = res.getUsernamesList().map(user => {
          return { username: user };
        });
        this.totalPageCount = res.getTotalPageSize();
      })
      .catch(err => {
        this.errHandler.handleError(err);
      });
  }
}
