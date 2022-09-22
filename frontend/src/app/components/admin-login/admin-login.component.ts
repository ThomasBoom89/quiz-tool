import {Component} from '@angular/core';
import {FormControl, FormGroup, Validators} from '@angular/forms';
import {AdminService} from '../../services/admin.service';

interface adminFormGroup {
  name: FormControl<string>;
  password: FormControl<string>;
}

@Component({
  selector: 'app-admin-login',
  templateUrl: './admin-login.component.html',
})
export class AdminLoginComponent {

  public adminForm: FormGroup<adminFormGroup>;

  constructor(private readonly adminService: AdminService) {
    this.adminForm = new FormGroup<adminFormGroup>({
      name: new FormControl<string>('', {nonNullable: true, validators: Validators.required}),
      password: new FormControl<string>('', {nonNullable: true, validators: Validators.required}),
    });
  }

  public onSubmitAdminForm() {
    this.adminService.login(this.adminForm.controls.name.value, this.adminForm.controls.password.value);
  }
}
