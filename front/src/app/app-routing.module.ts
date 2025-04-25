import { NgModule } from '@angular/core';
import { RouterModule, Routes } from '@angular/router';
import {LoginComponent} from './components/login/login.component';
import {QuestionPageComponent} from './components/question-page/question-page.component';

const routes: Routes = [
  { path: '', component: LoginComponent},
  { path: 'questions', component: QuestionPageComponent}
];

@NgModule({
  imports: [RouterModule.forRoot(routes)],
  exports: [RouterModule]
})
export class AppRoutingModule { }
