import {Component, Input} from '@angular/core';
import {PlayerState} from '../../enums/player-state';

@Component({
  selector: 'app-player-state',
  templateUrl: './player-state.component.html',
})
export class PlayerStateComponent {
  @Input() state!: PlayerState;

  constructor() {
  }
}
