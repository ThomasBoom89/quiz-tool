import {Component, OnInit} from '@angular/core';

@Component({
  selector: 'app-user-room',
  templateUrl: './user-room.component.html',
})
export class UserRoomComponent implements OnInit {
  currentQuestion: string = 'Hier könnte deine Frage stehen!';

  constructor() {
  }

  ngOnInit(): void {
  }

}
