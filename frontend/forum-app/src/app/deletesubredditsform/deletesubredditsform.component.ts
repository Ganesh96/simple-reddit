import { Component, OnInit } from '@angular/core';
import { FormBuilder, FormGroup, Validators } from '@angular/forms';
import { MatSnackBar } from '@angular/material/snack-bar';
import { SignupService } from '../signup.service';
import { Storage } from '../storage';

@Component({
  selector: 'app-deletesubredditsform',
  templateUrl: './deletesubredditsform.component.html',
  styleUrls: ['./deletesubredditsform.component.css']
})

export class DeletesubredditsformComponent implements OnInit {
  form: FormGroup = new FormGroup({});
  constructor(public signupService: SignupService, public snackBar: MatSnackBar, public fb: FormBuilder) {
    this.form = this.fb.group({
      username: ['', [Validators.required]],
      name: ['', [Validators.required]]
    })
   }
  
  get f() {
    return this.form.controls;
  }

  ngOnInit(): void {
  }

  getLoggedUsername() {
    if (Storage.isLoggedIn) {
      return Storage.username;
    }
    else {
      return "";
    }
  }

  deletesubreddit(username: string, name: string) {
    this.signupService.deletecommunity(username, name).subscribe((response: any) => {
      console.log(response);
    })

  }
}
