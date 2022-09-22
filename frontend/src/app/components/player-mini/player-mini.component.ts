import {Component, Input} from '@angular/core';
import {Player} from '../../interfaces/player';

@Component({
  selector: 'app-player-mini',
  templateUrl: './player-mini.component.html',
})
export class PlayerMiniComponent {
  @Input() player!: Player

  constructor() {
  }
}
