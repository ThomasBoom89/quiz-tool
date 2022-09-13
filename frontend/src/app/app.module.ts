import {NgModule} from '@angular/core';
import {BrowserModule} from '@angular/platform-browser';

import {AppRoutingModule} from './app-routing.module';
import {AppComponent} from './app.component';
import {UserComponent} from './pages/user/user.component';
import {AdminComponent} from './pages/admin/admin.component';
import {OverviewComponent} from './pages/overview/overview.component';
import {ReactiveFormsModule} from '@angular/forms';
import {UserRoomComponent} from './pages/user-room/user-room.component';
import {BuzzerComponent} from './components/buzzer/buzzer.component';
import {HttpClientModule} from '@angular/common/http';
import { StartComponent } from './pages/start/start.component';

@NgModule({
  declarations: [
    AppComponent,
    UserComponent,
    AdminComponent,
    OverviewComponent,
    UserRoomComponent,
    BuzzerComponent,
    StartComponent
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
