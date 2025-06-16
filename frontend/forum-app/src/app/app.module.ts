import { BrowserModule } from '@angular/platform-browser';
import { NgModule } from '@angular/core';

import { AppRoutingModule } from './app-routing.module';
import { AppComponent } from './app.component';
import { PostsComponent } from './posts/posts.component';
import { HttpClientModule } from '@angular/common/http';
import { BrowserAnimationsModule } from '@angular/platform-browser/animations';
import { MatCardModule } from '@angular/material/card';
import { MatButtonModule } from '@angular/material/button';
import { MatFormFieldModule } from '@angular/material/form-field';
import { MatInputModule } from '@angular/material/input';
import { MatSnackBarModule } from '@angular/material/snack-bar';
import { MatToolbarModule } from '@angular/material/toolbar';
import { MatIconModule } from '@angular/material/icon';
import { MatSidenavModule } from '@angular/material/sidenav';
import { MatListModule } from '@angular/material/list';
import { MatProgressSpinnerModule } from '@angular/material/progress-spinner';

import { NavbarComponent } from './navbar/navbar.component';
import { HomeComponent } from './home/home.component';
import { LoginComponent } from './login/login.component';
import { SignupformComponent } from './signupform/signupform.component';
import { SubredditsComponent } from './subreddits/subreddits.component';
import { NewsubredditsformComponent } from './newsubredditsform/newsubredditsform.component';
import { DeletesubredditsformComponent } from './deletesubredditsform/deletesubredditsform.component';
import { CommunitypageComponent } from './communitypage/communitypage.component';
import { NewpostformComponent } from './newpostform/newpostform.component';
import { PostpageComponent } from './postpage/postpage.component';
import { ProfileComponent } from './profile/profile.component';
import { DeleteuserformComponent } from './deleteuserform/deleteuserform.component';
import { EditpostComponent } from './editpost/editpost.component';
import { ContentpolicyComponent } from './contentpolicy/contentpolicy.component';
import { PrivacypolicyComponent } from './privacypolicy/privacypolicy.component';
import { ModpolicyComponent } from './modpolicy/modpolicy.component';
import { TermsandconditionsComponent } from './termsandconditions/termsandconditions.component';
import { ReactiveFormsModule } from '@angular/forms';
import { ApiService } from './api.service';

@NgModule({
  declarations: [
    AppComponent,
    PostsComponent,
    NavbarComponent,
    HomeComponent,
    LoginComponent,
    SignupformComponent,
    SubredditsComponent,
    NewsubredditsformComponent,
    DeletesubredditsformComponent,
    CommunitypageComponent,
    NewpostformComponent,
    PostpageComponent,
    ProfileComponent,
    DeleteuserformComponent,
    EditpostComponent,
    ContentpolicyComponent,
    PrivacypolicyComponent,
    ModpolicyComponent,
    TermsandconditionsComponent
  ],
  imports: [
    BrowserModule,
    AppRoutingModule,
    HttpClientModule,
    BrowserAnimationsModule,
    MatCardModule,
    MatButtonModule,
    MatFormFieldModule,
    MatInputModule,
    MatSnackBarModule,
    MatToolbarModule,
    MatIconModule,
    MatSidenavModule,
    MatListModule,
    MatProgressSpinnerModule,
    ReactiveFormsModule
  ],
  providers: [ApiService],
  bootstrap: [AppComponent]
})
export class AppModule { }