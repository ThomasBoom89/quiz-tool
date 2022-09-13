import {Injectable, OnDestroy} from '@angular/core';
import {webSocket, WebSocketSubject} from 'rxjs/webSocket';
import {Subscription} from 'rxjs';
import {WsEndpoint} from '../enums/ws-endpoint';

@Injectable({
  providedIn: 'root'
})
export class WebsocketService implements OnDestroy {
  private static readonly WS_ENPOINT_BASE = 'localhost:8898/';
  private static readonly WS_ENPOINT_USER = '/user/ws';
  private static readonly WS_ENPOINT_ADMIN = '/admin/ws';
  private static readonly WS_ENPOINT_OVERVIEW = '/overview/ws';

  private socket: WebSocketSubject<any> | undefined;
  private socketSubscription: Subscription | undefined;

  constructor() {
  }

  public ngOnDestroy() {
    this.socketShutdown();
  }

  public createWebsocketConnection(endpoint: WsEndpoint): boolean {
    // todo: handle token and validate it, else shutdown, return false
    this.socket = webSocket(this.getUrl(endpoint));
    this.socketSubscription = this.socket.subscribe((message: any) => {
      console.warn('message received: ', message);
    });

    return true
  }

  public sendMessage(message: string): void {
    this.socket?.next(message);
  }

  private socketShutdown() {
    this.socket?.complete();
    this.socketSubscription?.unsubscribe();
  }

  private getUrl(endpoint: WsEndpoint): string {
    switch (endpoint) {
      case WsEndpoint.user:
        return WebsocketService.WS_ENPOINT_BASE + WebsocketService.WS_ENPOINT_USER;
      case WsEndpoint.admin:
        return WebsocketService.WS_ENPOINT_BASE + WebsocketService.WS_ENPOINT_ADMIN;
      case WsEndpoint.overview:
        return WebsocketService.WS_ENPOINT_BASE + WebsocketService.WS_ENPOINT_OVERVIEW;
    }
  }
}
