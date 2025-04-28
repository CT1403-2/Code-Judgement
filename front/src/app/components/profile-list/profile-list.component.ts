import { Component } from '@angular/core';
import { Router } from '@angular/router';

@Component({
  selector: 'app-profile-list',
  standalone: false,
  templateUrl: './profile-list.component.html',
  styleUrl: './profile-list.component.css',
})
export class ProfileListComponent {
  currentUser!: string;

  constructor(private readonly router: Router) {}

  gotoMyProfile(tab: number) {
    if (tab == 1) {
      this.gotoProfile();
    }
  }

  gotoProfile(user: string = this.currentUser) {
    this.router.navigate(['profiles', user]);
  }
}
