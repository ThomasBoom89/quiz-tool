import {NgModule} from '@angular/core';
import {BrowserModule} from '@angular/platform-browser';

import {AppRoutingModule} from './app-routing.module';
import {AppComponent} from './app.component';
import {UserLoginComponent} from './components/user-login/user-login.component';
import {AdminComponent} from './pages/admin/admin.component';
import {OverviewComponent} from './pages/overview/overview.component';
import {ReactiveFormsModule} from '@angular/forms';
import {UserComponent} from './pages/user/user.component';
import {BuzzerComponent} from './components/buzzer/buzzer.component';
import {HttpClientModule} from '@angular/common/http';
import {StartComponent} from './pages/start/start.component';
import {AdminLoginComponent} from './components/admin-login/admin-login.component';
import {PlayerMiniComponent} from './components/player-mini/player-mini.component';
import {PlayerStateComponent} from './components/player-state/player-state.component';
import {PlayersComponent} from './components/players/players.component';

@NgModule({
  declarations: [
    AppComponent,
    UserLoginComponent,
    AdminComponent,
    OverviewComponent,
    UserComponent,
    BuzzerComponent,
    StartComponent,
    AdminLoginComponent,
    PlayerMiniComponent,
    PlayerStateComponent,
    PlayersComponent,
  ],
  imports: [
    BrowserModule,
    AppRoutingModule,
    ReactiveFormsModule,
    HttpClientModule,
  ],
  providers: [],
  bootstrap: [AppComponent]
})
export class AppModule {
}
