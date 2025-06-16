import { Component, OnInit } from '@angular/core';
import { FormBuilder, FormGroup, Validators } from '@angular/forms';
import { MatSnackBar } from '@angular/material/snack-bar';
import { Router } from '@angular/router';
import { ApiService } from '../api.service';

@Component({
  selector: 'app-signupform',
  templateUrl: './signupform.component.html',
  styleUrls: ['./signupform.component.css']
})
export class SignupformComponent implements OnInit {

  form: FormGroup;

  constructor(
    private fb: FormBuilder,
    private snackBar: MatSnackBar,
    private apiService: ApiService,
    private router: Router
  ) {
    this.form = this.fb.group({
      username: ['', Validators.required],
      password: ['', [
        Validators.required,
        Validators.minLength(8),
        Validators.pattern('^(?=.*[a-z])(?=.*[A-Z])(?=.*\\d).+$')
      ]]
    });
  }

  ngOnInit(): void {
  }

  get f() { return this.form.controls; }

  onSubmit() {
    if (this.form.valid) {
      this.apiService.post('user', this.form.value).subscribe((response: any) => {
        if (response.status == 201 && response.message == "success") {
          this.snackBar.open("User created successfully.", "Dismiss", { duration: 1500 });
          this.router.navigate(['/login']);
        } else {
          this.snackBar.open("Error creating user.", "Dismiss", { duration: 1500 });
        }
      });
    }
  }
}