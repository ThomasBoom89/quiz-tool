import {Component} from '@angular/core';
import {UserService} from '../../services/user.service';

@Component({
  selector: 'app-user',
  templateUrl: './user.component.html',
})
export class UserComponent {
  currentQuestion: string = 'Hier könnte deine Frage stehen!';

  constructor(
    public readonly userService: UserService,
  ) {
    this.userService.connect();
  }

  public buzzed() {
    this.userService.setBuzzed();
  }
}
