import {Injectable, OnDestroy} from '@angular/core';
import {webSocket, WebSocketSubject} from 'rxjs/webSocket';
import {Observable, Subject, Subscription} from 'rxjs';
import {WsEndpoint} from '../enums/ws-endpoint';
import {JwtService} from './jwt.service';
import {WebsocketEvent} from '../interfaces/websocket-event';

@Injectable({
  providedIn: 'root'
})
export class WebsocketService implements OnDestroy {
  private static readonly WS_ENPOINT_BASE = 'ws://localhost:8898/api/v1';
  private static readonly WS_ENPOINT_USER = '/user/ws';
  private static readonly WS_ENPOINT_ADMIN = '/admin/ws';
  private static readonly WS_ENPOINT_OVERVIEW = '/overview/ws';

  private socket: WebSocketSubject<any> | undefined;
  private socketSubscription: Subscription | undefined;

  private websocketSubject: Subject<WebsocketEvent> = new Subject<WebsocketEvent>();
  public readonly websocketObservable: Observable<WebsocketEvent> = this.websocketSubject.asObservable();

  constructor(private readonly jwtService: JwtService) {
  }

  public ngOnDestroy() {
    this.socketShutdown();
  }

  public createWebsocketConnection(endpoint: WsEndpoint): boolean {
    this.socket = webSocket(this.getUrl(endpoint));
    this.socketSubscription = this.socket.subscribe((websocketEvent: WebsocketEvent) => {
      // todo validate event
      this.websocketSubject.next(websocketEvent);
    });

    return true
  }

  public sendMessage(message: any): void {
    this.socket?.next(message);
  }

  private socketShutdown() {
    this.socket?.complete();
    this.socketSubscription?.unsubscribe();
  }

  private getUrl(endpoint: WsEndpoint): string {
    let url = '';
    switch (endpoint) {
      case WsEndpoint.user:
        url = WebsocketService.WS_ENPOINT_BASE + WebsocketService.WS_ENPOINT_USER;
        break;
      case WsEndpoint.admin:
        url = WebsocketService.WS_ENPOINT_BASE + WebsocketService.WS_ENPOINT_ADMIN;
        break;
      case WsEndpoint.overview:
        url = WebsocketService.WS_ENPOINT_BASE + WebsocketService.WS_ENPOINT_OVERVIEW;
        break;
    }

    url = url + '?token=' + (this.jwtService.getToken());
    console.warn('url: ', url);

    return url;
  }
}
