import {Injectable, OnDestroy} from '@angular/core';
import {WebsocketService} from './websocket.service';
import {WebsocketEvent} from '../interfaces/websocket-event';
import {Subscription} from 'rxjs';

@Injectable({
  providedIn: 'root'
})
export class QuizService implements OnDestroy {
  private websocketServiceSubscription: Subscription;

  constructor(private readonly websocketService: WebsocketService) {
    this.websocketServiceSubscription = this.websocketService.websocketObservable
      .subscribe((websocketEvent: WebsocketEvent) => {
        console.warn('event received: ', websocketEvent);
      });
  }

  ngOnDestroy() {
    this.websocketServiceSubscription.unsubscribe();
  }
}
