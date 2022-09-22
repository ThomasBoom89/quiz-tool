import {Component} from '@angular/core';
import {Player} from '../../interfaces/player';
import {PlayerState} from '../../enums/player-state';
import {FormControl, FormGroup, Validators} from '@angular/forms';

interface questionFormGroup {
  question: FormControl<string>;
}

@Component({
  selector: 'app-admin',
  templateUrl: './admin.component.html',
})
export class AdminComponent {

  public players: Player[] = [
    {id: 1, name: 'test', points: 4, state: PlayerState.none},
    {id: 2, name: 'auch', points: 7, state: PlayerState.none},
    {id: 2, name: 'noch', points: 7, state: PlayerState.none},
    {id: 1, name: 'test', points: 4, state: PlayerState.none},
    {id: 2, name: 'auch', points: 7, state: PlayerState.none},
    {id: 2, name: 'diesisteinwirklichlangername', points: 7, state: PlayerState.none},
    {id: 1, name: 'test', points: 4, state: PlayerState.none},
    {id: 2, name: 'auch', points: 7, state: PlayerState.none},
    {id: 2, name: 'noch', points: 7, state: PlayerState.active},
    {id: 1, name: 'test', points: 4, state: PlayerState.none},
    {id: 2, name: 'auch', points: 7, state: PlayerState.none},
    {id: 2, name: 'noch', points: 7, state: PlayerState.none},
    {id: 1, name: 'test', points: 4, state: PlayerState.none},
    {id: 2, name: 'auch', points: 7, state: PlayerState.blocked},
    {id: 2, name: 'noch', points: 7, state: PlayerState.none},
    {id: 1, name: 'test', points: 4, state: PlayerState.none},
    {id: 2, name: 'auch', points: 7, state: PlayerState.none},
    {id: 2, name: 'noch', points: 7, state: PlayerState.none},
  ];

  public questionForm: FormGroup<questionFormGroup>;

  constructor() {
    this.questionForm = new FormGroup<questionFormGroup>({
      question: new FormControl<string>('', {nonNullable: true, validators: Validators.required}),
    });
  }

  public onSubmitQuestionForm() {
  }
}
