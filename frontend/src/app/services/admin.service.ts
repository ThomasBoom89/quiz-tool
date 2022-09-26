import {Injectable} from '@angular/core';
import {Router} from '@angular/router';

@Injectable({
  providedIn: 'root'
})
export class AdminService {

  private static readonly LOGIN_URL = 'api/v1/admin/login';

  constructor(private readonly router: Router) {
  }

  public login(name: string, password: string) {
    console.warn('name: ', name, 'password: ', password);

    this.router.navigateByUrl('/admin');
  }
}
