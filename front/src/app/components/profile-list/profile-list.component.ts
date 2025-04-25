import { Component } from '@angular/core';

@Component({
  selector: 'app-profile-list',
  standalone: false,
  templateUrl: './profile-list.component.html',
  styleUrl: './profile-list.component.css'
})
export class ProfileListComponent {
  myProfile(tab: number) {
    if (tab == 0) return;
  }
}
