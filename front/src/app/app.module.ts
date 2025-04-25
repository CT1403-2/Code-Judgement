import { NgModule } from '@angular/core';
import { BrowserModule } from '@angular/platform-browser';
import {FormsModule} from '@angular/forms';
import { NgChartjsModule } from 'ng-chartjs';

import { AppRoutingModule } from './app-routing.module';
import { AppComponent } from './app.component';
import { LoginComponent } from './components/login/login.component';
import { TabManagerComponent } from './components/tab-manager/tab-manager.component';
import { QuestionPageComponent } from './components/question-page/question-page.component';
import { QuestionListComponent } from './components/question-page/question-list/question-list.component';
import { NewQuestionComponent } from './components/question-page/new-question/new-question.component';
import { ProfileListComponent } from './components/profile-list/profile-list.component';
import { ProfileComponent } from './components/profile/profile.component';
import { QuestionComponent } from './components/question/question.component';
import { SubmitComponent } from './components/question/submit/submit.component';
import { ProfileDetailComponent } from './components/profile/profile-detail/profile-detail.component';
import { QuestionDetailComponent } from './components/question/question-detail/question-detail.component';
import { SubmissionListComponent } from './components/submission-list/submission-list.component';
import { TableComponent } from './components/table/table.component';

@NgModule({
  declarations: [
    AppComponent,
    LoginComponent,
    TabManagerComponent,
    QuestionPageComponent,
    QuestionListComponent,
    NewQuestionComponent,
    ProfileListComponent,
    ProfileComponent,
    QuestionComponent,
    SubmitComponent,
    ProfileDetailComponent,
    QuestionDetailComponent,
    SubmissionListComponent,
    TableComponent
  ],
  imports: [
    BrowserModule,
    FormsModule,
    NgChartjsModule,
    AppRoutingModule
  ],
  providers: [],
  bootstrap: [AppComponent]
})
export class AppModule { }
