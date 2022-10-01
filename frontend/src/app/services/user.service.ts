import {Injectable} from '@angular/core';
import {WebsocketService} from './websocket.service';
import {Router} from '@angular/router';
import {Observable} from 'rxjs';
import {HttpClient} from '@angular/common/http';
import {JwtService} from './jwt.service';
import {WsEndpoint} from '../enums/ws-endpoint';
import {QuizService} from './quiz.service';
import {UserAction} from '../enums/user-action';
import {LoginResponse} from '../interfaces/login-response';
import {Player} from '../interfaces/player';

@Injectable({
  providedIn: 'root'
})
export class UserService {

  private static readonly LOGIN_URL = 'api/v1/user/login';
  private id: string = '';

  constructor(
    public quizService: QuizService,
    private readonly websocketService: WebsocketService,
    private readonly router: Router,
    private readonly httpClient: HttpClient,
    private readonly jwtService: JwtService,
  ) {
  }

  public connect(): void {
    this.websocketService.createWebsocketConnection(WsEndpoint.user);
  }

  public setBuzzed(): void {
    this.websocketService.sendMessage({action: UserAction.buzzed, payload: this.id});
  }

  public register(name: string, roomId: string): void {
    console.warn('try to register: ', name, roomId);
    this.getToken(name, roomId).subscribe((response: LoginResponse) => {
      console.warn('received token: ', response.token);
      this.jwtService.setToken(response.token);
      this.id = response.id;
      this.router.navigateByUrl('user/room/' + roomId);
    });
  }

  private getToken(name: string, roomId: string): Observable<LoginResponse> {
    const param = {
      name: name,
      roomId: roomId,
    };
    return this.httpClient.post<LoginResponse>(UserService.LOGIN_URL, param);
  }

  public getPlayers(): Player[] {
    return this.quizService.getPlayers().filter((player: Player) => {
      return player.id !== this.id;
    });
  }

  public getPoints(): number {
    const player = this.quizService.getPlayer(this.id);
    if (player === undefined) {
      return 0;
    }

    return player.points;
  }

  public getState(): number {
    const player = this.quizService.getPlayer(this.id);
    if (player === undefined) {
      return 0;
    }

    return player.state;
  }
}
