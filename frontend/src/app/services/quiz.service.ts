import {Injectable, OnDestroy} from '@angular/core';
import {WebsocketService} from './websocket.service';
import {WebsocketEvent} from '../interfaces/websocket-event';
import {Subscription} from 'rxjs';
import {QuizAction} from '../enums/quiz-action';
import {Player} from '../interfaces/player';

@Injectable({
  providedIn: 'root'
})
export class QuizService implements OnDestroy {
  private playerMap: Map<string, Player> = new Map;
  private websocketServiceSubscription: Subscription;
  private question: string = '';

  constructor(private readonly websocketService: WebsocketService) {
    this.websocketServiceSubscription = this.websocketService.websocketObservable
      .subscribe((websocketEvent: WebsocketEvent) => {
        console.warn('event received: ', websocketEvent);
        switch (websocketEvent.action) {
          case QuizAction.Init:
            if (websocketEvent.payload === null) {
              this.playerMap.clear();
              break;
            }
            const playerMapClone = new Map<string, Player>(this.playerMap);
            for (const player of websocketEvent.payload as Player[]) {
              if (playerMapClone.has(player.id)) {
                playerMapClone.delete(player.id);
              }
              this.addPlayer(player);
            }
            playerMapClone.forEach((player: Player) => {
              this.removePlayer(player.id);
            });
            break;
          case QuizAction.UserEntered:
            this.addPlayer(websocketEvent.payload as Player);
            break;
          case QuizAction.UserLeft:
            this.removePlayer(websocketEvent.payload as string);
            break;
          case QuizAction.UserStatus:
            this.updatePlayer(websocketEvent.payload as Player);
            break;
          case QuizAction.NewRound:
            this.question = websocketEvent.payload;
            break;
        }
      });
  }

  ngOnDestroy() {
    this.websocketServiceSubscription.unsubscribe();
  }

  public getPlayers(): Player[] {
    return [...this.playerMap.values()]
  }

  public getQuestion(): string {
    return this.question;
  }

  private addPlayer(player: Player): void {
    this.playerMap.set(player.id, player);
  }

  private updatePlayer(player: Player): void {
    if (this.playerMap.has(player.id) === undefined) {
      return;
    }

    this.playerMap.set(player.id, player);
  }

  private removePlayer(id: string): void {
    this.playerMap.delete(id);
  }
}
