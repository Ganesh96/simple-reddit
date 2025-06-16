import { Component, OnInit } from '@angular/core';
import { FormBuilder, FormGroup, Validators } from '@angular/forms';
import { MatSnackBar } from '@angular/material/snack-bar';
import { Router } from '@angular/router';
import { ApiService } from '../api.service';
import { Storage } from '../storage';

@Component({
  selector: 'app-newpostform',
  templateUrl: './newpostform.component.html',
  styleUrls: ['./newpostform.component.css']
})
export class NewpostformComponent implements OnInit {

  form: FormGroup;

  constructor(
    private fb: FormBuilder,
    private snackBar: MatSnackBar,
    private apiService: ApiService,
    private router: Router
  ) {
    this.form = this.fb.group({
      title: ['', Validators.required],
      community: ['', Validators.required],
      body: ['', Validators.required]
    });
  }

  ngOnInit(): void {
  }

  onSubmit() {
    if (this.form.valid) {
      const { title, community, body } = this.form.value;
      this.apiService.post('post', {
        username: Storage.username,
        community,
        title,
        body
      }).subscribe((response: any) => {
        if (response.status == 201 && response.message == "success") {
          this.snackBar.open("New post created.", "Dismiss", { duration: 1500 });
          this.router.navigate(['/post', response.data.post_id]);
        } else {
          this.snackBar.open("Error creating post.", "Dismiss", { duration: 1500 });
        }
      });
    }
  }
}