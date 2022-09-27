import {Injectable} from '@angular/core';
import {Router} from '@angular/router';
import {JwtService} from './jwt.service';
import {HttpClient} from '@angular/common/http';
import {LoginResponse} from '../interfaces/login-response';
import {WsEndpoint} from '../enums/ws-endpoint';
import {WebsocketService} from './websocket.service';
import {AdminAction} from '../enums/admin-action';

@Injectable({
  providedIn: 'root'
})
export class AdminService {

  private static readonly LOGIN_URL = 'api/v1/admin/login';

  constructor(
    private readonly router: Router,
    private readonly http: HttpClient,
    private readonly jwtService: JwtService,
    private readonly websocketService: WebsocketService,
  ) {
  }

  public connect(): void {
    this.websocketService.createWebsocketConnection(WsEndpoint.admin);
  }

  public setAction(action: AdminAction, payload: string = ''): void {
    this.websocketService.sendMessage({action: action, payload: payload})
  }

  public login(name: string, password: string) {
    const param = {
      name,
      password,
    };
    this.http.post<LoginResponse>(AdminService.LOGIN_URL, param)
      .subscribe((response: LoginResponse) => {
        console.warn('received token: ', response.token);
        this.jwtService.setToken(response.token);
        // todo check if token was transmitted or an error occurred
        this.router.navigateByUrl('/admin');
      });

  }
}
