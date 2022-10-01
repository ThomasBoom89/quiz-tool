import {Component, Input, OnInit} from '@angular/core';
import {UserService} from '../../services/user.service';
import {AdminAction} from '../../enums/admin-action';
import {AdminService} from '../../services/admin.service';
import {Player} from '../../interfaces/player';

@Component({
  selector: 'app-players',
  templateUrl: './players.component.html',
})
export class PlayersComponent implements OnInit {

  @Input() players!: Player[];
  @Input() isAdmin: boolean = false;

  constructor(
    private readonly adminService: AdminService,
  ) {
  }

  ngOnInit(): void {
  }


  public removePlayer(playerId: string) {
    this.adminService.setAction(AdminAction.RemoveUser, playerId);
  }
}
