import { Component, OnInit } from '@angular/core';
import { FormBuilder, FormGroup, Validators } from '@angular/forms';
import { MatSnackBar } from '@angular/material/snack-bar';
import { Router } from '@angular/router';
import { ApiService } from '../api.service';

@Component({
  selector: 'app-newsubredditsform',
  templateUrl: './newsubredditsform.component.html',
  styleUrls: ['./newsubredditsform.component.css']
})
export class NewsubredditsformComponent implements OnInit {

  form: FormGroup;

  constructor(
    private fb: FormBuilder,
    private snackBar: MatSnackBar,
    private apiService: ApiService,
    private router: Router
  ) {
    this.form = this.fb.group({
      name: ['', Validators.required]
    });
  }

  ngOnInit(): void {
  }

  onSubmit() {
    if (this.form.valid) {
      this.apiService.post('community', this.form.value).subscribe((response: any) => {
        if (response.status == 201 && response.message == "success") {
          this.snackBar.open("New community created.", "Dismiss", { duration: 1500 });
          this.router.navigate(['/r', this.form.value.name]);
        } else {
          this.snackBar.open("Error creating community.", "Dismiss", { duration: 1500 });
        }
      });
    }
  }
}