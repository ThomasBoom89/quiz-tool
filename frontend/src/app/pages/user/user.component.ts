import {Component, OnInit} from '@angular/core';

@Component({
  selector: 'app-user',
  templateUrl: './user.component.html',
})
export class UserComponent implements OnInit {
  currentQuestion: string = 'Hier könnte deine Frage stehen!';

  constructor() {
  }

  ngOnInit(): void {
  }

  public buzzed() {
    console.warn('es wurde gebuzzed');
  }
}
