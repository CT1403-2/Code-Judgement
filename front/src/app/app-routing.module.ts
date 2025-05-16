import { NgModule } from '@angular/core';
import { RouterModule, Routes } from '@angular/router';
import { LoginComponent } from './components/login/login.component';
import { QuestionPageComponent } from './components/question-page/question-page.component';
import { QuestionComponent } from './components/question/question.component';
import { ProfileListComponent } from './components/profile-list/profile-list.component';
import { ProfileComponent } from './components/profile/profile.component';
import { ErrorComponent } from './components/error/error.component';

const routes: Routes = [
  { path: '', component: LoginComponent },
  { path: 'questions', component: QuestionPageComponent },
  { path: 'questions/:id', component: QuestionComponent },
  { path: 'profiles', component: ProfileListComponent },
  { path: 'profiles/:id', component: ProfileComponent },
  { path: 'error/:id', component: ErrorComponent },
  { path: '**', redirectTo: '/error/404' }
];

@NgModule({
  imports: [RouterModule.forRoot(routes)],
  exports: [RouterModule]
})
export class AppRoutingModule {}
