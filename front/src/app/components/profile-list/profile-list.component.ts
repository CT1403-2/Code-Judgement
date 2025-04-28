import { Component, OnInit } from '@angular/core';
import { Router } from '@angular/router';
import { ManagerService } from '../../services/manager.service';

@Component({
  selector: 'app-profile-list',
  standalone: false,
  templateUrl: './profile-list.component.html',
  styleUrl: './profile-list.component.css',
})
export class ProfileListComponent implements OnInit {
  currentUser!: string;
  profiles!: { username: string }[];

  constructor(
    private readonly router: Router,
    private readonly manager: ManagerService,
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
      .GetProfile({
        value: '',
      })
      .then((res) => {
        this.currentUser = res.username;
      })
      .catch((err) => {});
    this.manager
      .GetProfiles({
        filters: [],
      })
      .then((res) => {
        this.profiles = res.usernames.map((user) => {
          return { username: user };
        });
      })
      .catch((err) => {});
  }
}
