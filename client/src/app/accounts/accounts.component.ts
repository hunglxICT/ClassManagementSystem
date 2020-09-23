import { Component } from '@angular/core';
import { first } from 'rxjs/operators';
import { Location } from '@angular/common';
import { ActivatedRoute } from '@angular/router';

import { User, UserBackend } from '@app/_models';
import { UserService } from '@app/_services';

@Component({
  selector: 'app-accounts',
  templateUrl: './accounts.component.html',
  styleUrls: ['./accounts.component.css']
})
export class AccountsComponent {
    loading = false;
    users: UserBackend[];

    constructor(
      private userService: UserService,
      private route: ActivatedRoute,
      private location: Location
    ) { }

    ngOnInit() {
        this.loading = true;
        this.userService.getAll().pipe(first()).subscribe(users => {
            this.loading = false;
            this.users = users;
        });
    }
    
    goBack(): void {
      this.location.back();
    }
    
    delete(account: UserBackend): void {
      this.userService.deleteByID(account['Id'])
        .subscribe();
      this.userService.getAll().pipe(first()).subscribe(users => {
            this.loading = false;
            this.users = users;
        });
    }
}
