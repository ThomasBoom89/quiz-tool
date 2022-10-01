import {Component, EventEmitter, Input, Output} from '@angular/core';
import {Player} from '../../interfaces/player';

@Component({
  selector: 'app-player-mini',
  templateUrl: './player-mini.component.html',
})
export class PlayerMiniComponent {
  @Input() player!: Player
  @Input() isAdmin: boolean = false;

  @Output() removePlayer: EventEmitter<string> = new EventEmitter<string>();

  constructor() {
  }
}
