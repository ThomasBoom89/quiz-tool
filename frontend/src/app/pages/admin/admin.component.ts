import {Component} from '@angular/core';
import {Player} from '../../interfaces/player';
import {PlayerState} from '../../enums/player-state';
import {FormControl, FormGroup, Validators} from '@angular/forms';
import {AdminService} from '../../services/admin.service';
import {AdminAction} from '../../enums/admin-action';
import {QuizService} from '../../services/quiz.service';

interface questionFormGroup {
  question: FormControl<string>;
}

@Component({
  selector: 'app-admin',
  templateUrl: './admin.component.html',
})
export class AdminComponent {

  public players: Player[] = [
    {id: '1', name: 'test', points: 4, state: PlayerState.none},
    {id: '2', name: 'auch', points: 7, state: PlayerState.none},
    {id: '2', name: 'noch', points: 7, state: PlayerState.none},
    {id: '1', name: 'test', points: 4, state: PlayerState.none},
    {id: '2', name: 'auch', points: 7, state: PlayerState.none},
    {id: '2', name: 'diesisteinwirklichlangername', points: 7, state: PlayerState.none},
    {id: '1', name: 'test', points: 4, state: PlayerState.none},
    {id: '2', name: 'auch', points: 7, state: PlayerState.none},
    {id: '2', name: 'noch', points: 7, state: PlayerState.active},
    {id: '1', name: 'test', points: 4, state: PlayerState.none},
    {id: '2', name: 'auch', points: 7, state: PlayerState.none},
    {id: '2', name: 'noch', points: 7, state: PlayerState.none},
    {id: '1', name: 'test', points: 4, state: PlayerState.none},
    {id: '2', name: 'auch', points: 7, state: PlayerState.blocked},
    {id: '2', name: 'noch', points: 7, state: PlayerState.none},
    {id: '1', name: 'test', points: 4, state: PlayerState.none},
    {id: '2', name: 'auch', points: 7, state: PlayerState.none},
    {id: '2', name: 'noch', points: 7, state: PlayerState.none},
  ];

  public questionForm: FormGroup<questionFormGroup>;

  constructor(private readonly adminService: AdminService, public readonly quizService: QuizService) {
    this.adminService.connect();
    this.questionForm = new FormGroup<questionFormGroup>({
      question: new FormControl<string>('', {nonNullable: true, validators: Validators.required}),
    });
  }

  public onSubmitQuestionForm() {
    console.warn(this.questionForm.controls);
    this.adminService.setAction(AdminAction.StartNewQuestion, this.questionForm.controls.question.value);
  }

  public setCorrectAnswer() {
    this.adminService.setAction(AdminAction.SetCorrectAnswer);
  }

  public setWrongAnswer() {
    this.adminService.setAction(AdminAction.SetWrongAnswer);
  }
}
