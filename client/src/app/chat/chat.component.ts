import { Component, OnInit } from '@angular/core';
import { ActivatedRoute } from '@angular/router';
import { Location } from '@angular/common';
import { first } from 'rxjs/operators';


import { User, UserBackend, Chat } from '@app/_models';
import { UserService, ChatService } from '@app/_services';


@Component({
  selector: 'app-chat',
  templateUrl: './chat.component.html',
  styleUrls: ['./chat.component.less']
})
export class ChatComponent implements OnInit {
  
  chat: Chat[];
  newmessage: Chat;
  
  constructor(
    private route: ActivatedRoute,
    private chatService: ChatService,
    private accountService: UserService,
    private location: Location
  ) { }

  ngOnInit() { 
    this.newmessage = new Chat;
    this.getChat();
  }
  
  getChat(): void {
    const id = +this.route.snapshot.paramMap.get('id');
    this.chatService.getChat(id).subscribe(chat => this.chat = chat['result']);
  }
  
  sendMessage(): void {
    const id = +this.route.snapshot.paramMap.get('id');
    this.newmessage.Receiverid = id;
    this.chatService.sendMessage(this.newmessage).subscribe(_ => {this.getChat(); this.newmessage = new Chat;});
  }
  
  deleteMessage(message: Chat): void {
    this.chatService.deleteMessage(message).subscribe(_ => {this.getChat()});
  }
  
  editMessage(message: Chat): void {
    this.chatService.editMessage(message).subscribe(_ => {this.getChat()});
  }
  
}
