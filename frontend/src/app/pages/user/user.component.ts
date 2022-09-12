import {Component} from '@angular/core';
import {FormControl, FormGroup, Validators} from '@angular/forms';
import {UserService} from '../../services/user.service';

interface userFormGroup {
  name: FormControl<string>;
  roomId: FormControl<string>;
}

@Component({
  selector: 'app-user',
  templateUrl: './user.component.html',
})
export class UserComponent {
  public userForm: FormGroup<userFormGroup>;

  constructor(private readonly userService: UserService) {
    this.userForm = new FormGroup<userFormGroup>({
      name: new FormControl<string>('', {nonNullable: true, validators: Validators.required}),
      roomId: new FormControl<string>('', {nonNullable: true, validators: Validators.required}),
    });
  }

  public onSubmitUserForm() {
    this.userService.register(this.userForm.controls.name.value, this.userForm.controls.roomId.value);
  }
}
