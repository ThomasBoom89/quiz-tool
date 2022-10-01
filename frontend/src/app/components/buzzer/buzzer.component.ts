import {Component, EventEmitter, HostListener, Input, Output} from '@angular/core';

@Component({
  selector: 'app-buzzer',
  templateUrl: './buzzer.component.html',
})
export class BuzzerComponent {
  @Input() isBuzzed: boolean = false;

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
