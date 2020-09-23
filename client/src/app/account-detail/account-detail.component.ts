import { Component, OnInit } from '@angular/core';
import { ActivatedRoute } from '@angular/router';
import { Location } from '@angular/common';
import { first } from 'rxjs/operators';

import { User, UserBackend, Chat } from '@app/_models';
import { UserService, ChatService } from '@app/_services';

@Component({
  selector: 'app-account-detail',
  templateUrl: './account-detail.component.html',
  //styleUrls: ['./account-detail.component.css']
})
export class AccountDetailComponent implements OnInit {

  account: UserBackend;

  constructor(
    private route: ActivatedRoute,
    private accountService: UserService,
    private chatService: ChatService,
    private location: Location
  ) {}

  ngOnInit(): void {
    this.getAccount();
    
  }

  getAccount(): void {
    const id = +this.route.snapshot.paramMap.get('id');
    this.accountService.getByID(id)
      .subscribe(acc => this.account = acc['result']);
  }

  goBack(): void {
    this.location.back();
  }

  save(): void {
  //alert(JSON.stringify(this.account['result']));
    this.accountService.saveByID(this.account)
      .subscribe(() => this.goBack());
  }
  
  delete(): void {
    this.accountService.deleteByID(this.account['Id'])
      .subscribe(() => this.goBack());
  }

}
