import {PlayerState} from '../enums/player-state';

export interface Player {
  id: number;
  name: string;
  points: number;
  state: PlayerState;
}
