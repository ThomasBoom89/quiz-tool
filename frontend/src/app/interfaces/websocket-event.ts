import {QuizAction} from '../enums/quiz-action';

export interface WebsocketEvent {
  action: QuizAction;
  payload: any;
}
