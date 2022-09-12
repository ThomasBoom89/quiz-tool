import {Injectable} from '@angular/core';
import {WebsocketService} from './websocket.service';
import {Router} from '@angular/router';
import {Observable} from 'rxjs';
import {HttpClient} from '@angular/common/http';

interface TokenResponse {
  token: string;
}

@Injectable({
  providedIn: 'root'
})
export class UserService {

  private static readonly LOGIN_URL = 'api/v1/user/login';

  constructor(
    private readonly websocketService: WebsocketService,
    private readonly router: Router,
    private readonly httpClient: HttpClient,
  ) {
  }

  public register(name: string, roomId: string): void {
    console.warn('try to register: ', name, roomId);
    this.getToken(name, roomId).subscribe((response: TokenResponse) => {
      console.warn('received token: ', response.token);
      // todo: check token and save it
      this.router.navigateByUrl('user/room/' + roomId);
    });
  }

  public getToken(name: string, roomId: string): Observable<TokenResponse> {
    const param = {
      name: name,
      roomId: roomId,
    };
    return this.httpClient.post<TokenResponse>(UserService.LOGIN_URL, param);
  }
}
