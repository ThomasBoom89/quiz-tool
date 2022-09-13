import {Component, EventEmitter, HostListener, Output} from '@angular/core';

@Component({
  selector: 'app-buzzer',
  templateUrl: './buzzer.component.html',
})
export class BuzzerComponent {
  @Output() buzzed: EventEmitter<boolean> = new EventEmitter<boolean>();

  @HostListener('document:keypress', ['$event'])
  keydown(e: KeyboardEvent) {
    if (e.key === ' ') {
      this.clicked()
    }
  }

  public clicked() {
    this.buzzed.emit(true);
  }
}
