import {Component} from '@angular/core';
import {UserService} from '../../services/user.service';
import {QuizService} from '../../services/quiz.service';
import {PlayerState} from '../../enums/player-state';

@Component({
  selector: 'app-user',
  templateUrl: './user.component.html',
})
export class UserComponent {
  constructor(
    public readonly userService: UserService,
    public readonly quizService: QuizService,
  ) {
    this.userService.connect();
  }

  public buzzed() {
    this.userService.setBuzzed();
  }

  public isBuzzed(): boolean {
    return this.userService.getState() === PlayerState.active;
  }
}
