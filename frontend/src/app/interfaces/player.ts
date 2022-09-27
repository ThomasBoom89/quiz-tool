import {PlayerState} from '../enums/player-state';

export interface Player {
  id: string;
  name: string;
  points: number;
  state: PlayerState;
}
