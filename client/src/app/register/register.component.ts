import { Component, OnInit } from '@angular/core';
import { ActivatedRoute } from '@angular/router';
import { Location } from '@angular/common';

import { UserBackend } from '@app/_models';
import { UserService } from '@app/_services';

@Component({
  selector: 'app-register',
  templateUrl: './register.component.html',
  //styleUrls: ['./register.component.css']
})
export class RegisterComponent implements OnInit {
  
  account: UserBackend;
  
  constructor(
    private route: ActivatedRoute,
    private accountService: UserService,
    private location: Location
  ) { }

  ngOnInit() {
    this.account = new UserBackend;
    //alert(JSON.stringify(this.account));
  }
  
  goBack(): void {
    this.location.back();
  }
  
  save(): void {
    this.accountService.register(this.account)
      .subscribe(() => this.goBack());
  }
}
