import {Component} from '@angular/core';
import {QuizService} from '../../services/quiz.service';

@Component({
  selector: 'app-current-question',
  templateUrl: './current-question.component.html',
})
export class CurrentQuestionComponent {
  constructor(public readonly quizService: QuizService) {
  }

}
