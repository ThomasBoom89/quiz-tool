import {Injectable} from '@angular/core';
import {WsEndpoint} from '../enums/ws-endpoint';
import {WebsocketService} from './websocket.service';

@Injectable({
  providedIn: 'root'
})
export class OverviewService {

  constructor(private readonly websocketService: WebsocketService) {
  }

  public connect(): void {
    this.websocketService.createWebsocketConnection(WsEndpoint.overview);
  }
}
