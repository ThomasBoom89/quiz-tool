import {Component} from '@angular/core';
import {QuizService} from '../../services/quiz.service';
import {OverviewService} from '../../services/overview.service';

@Component({
  selector: 'app-overview',
  templateUrl: './overview.component.html',
})
export class OverviewComponent {

  constructor(
    public readonly overviewService: OverviewService,
    public readonly quizService: QuizService,
  ) {
    this.overviewService.connect();
  }
}
