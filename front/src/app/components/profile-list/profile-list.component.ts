import { Component } from '@angular/core';
import {Router} from '@angular/router';

@Component({
  selector: 'app-profile-list',
  standalone: false,
  templateUrl: './profile-list.component.html',
  styleUrl: './profile-list.component.css'
})
export class ProfileListComponent {
  currentUser: string = "ahshqhir";
  activeTab: number = 0;

  constructor(private readonly router: Router) {}

  myProfile(tab: number) {
    if (tab == 1) {
      this.profile();
    }
  }

  profile(user: string = this.currentUser) {
    this.router.navigate(['profiles', user]);
  }
}
