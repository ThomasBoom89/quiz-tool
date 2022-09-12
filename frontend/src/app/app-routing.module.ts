import {NgModule} from '@angular/core';
import {RouterModule, Routes} from '@angular/router';
import {UserComponent} from './pages/user/user.component';
import {AdminComponent} from './pages/admin/admin.component';
import {OverviewComponent} from './pages/overview/overview.component';
import {UserRoomComponent} from './pages/user-room/user-room.component';

const routes: Routes = [
  {
    path: 'user',
    component: UserComponent,
  },
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
  {path: '**', redirectTo: '/'},
];

@NgModule({
  imports: [RouterModule.forRoot(routes)],
  exports: [RouterModule]
})


export class AppRoutingModule {
}
