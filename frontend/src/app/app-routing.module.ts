import {NgModule} from '@angular/core';
import {RouterModule, Routes} from '@angular/router';
import {AdminComponent} from './pages/admin/admin.component';
import {OverviewComponent} from './pages/overview/overview.component';
import {UserRoomComponent} from './pages/user-room/user-room.component';
import {StartComponent} from './pages/start/start.component';

const routes: Routes = [
  {
    path: 'user/room/:roomId',
    component: UserRoomComponent,
  },
  {
    path: 'admin',
    component: AdminComponent,
  },
  {
    path: 'overview',
    component: OverviewComponent,
  },
  {path: '**', component: StartComponent},
];

@NgModule({
  imports: [RouterModule.forRoot(routes)],
  exports: [RouterModule]
})


export class AppRoutingModule {
}
