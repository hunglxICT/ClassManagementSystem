import { Injectable } from '@angular/core';
import { HttpClient, HttpHeaders } from '@angular/common/http';

import { environment } from '@environments/environment';
import { User, UserBackend, Chat } from '@app/_models';
import { catchError, map, tap } from 'rxjs/operators';
import { Observable, of } from 'rxjs';

@Injectable({ providedIn: 'root' })
export class ChatService {
    constructor(private http: HttpClient) { }
    httpOptions = {
    	headers: new HttpHeaders({ 'Content-Type': 'application/json' })
    };
    
    getChat(id: number) {
        const url = `${environment.apiUrl}/get-messenger/${id}`;
    	return this.http.get<Chat[]>(url).pipe(
      		tap(_ => this.log(`fetched messenger for id=${id}`)),
      		catchError(this.handleError<Chat[]>(`getChat id=${id}`))
    	);
    }
    
    sendMessage(newmessage: Chat) {
        const url = `${environment.apiUrl}/send-messenger`;
    	return this.http.post<Chat>(url, newmessage, this.httpOptions).pipe(
    		tap(_ => this.log(`new message to user id=${newmessage.Receiverid}`)),
    		catchError(this.handleError<any>('sendMessage'))
    	);
    }
    
    deleteMessage(message: Chat) {
        const url = `${environment.apiUrl}/delete-messenger`;
        return this.http.post(url, message, this.httpOptions).pipe(
    		tap(_ => this.log(`delete message id=${message.Receiverid}`)),
    		catchError(this.handleError<any>('deleteMessage'))
    	);
    }
    
    editMessage(message: Chat) {
        const url = `${environment.apiUrl}/edit-messenger`;
        return this.http.post(url, message, this.httpOptions).pipe(
    		tap(_ => this.log(`delete message id=${message.Receiverid}`)),
    		catchError(this.handleError<any>('deleteMessage'))
    	);
    }
    
    private handleError<T>(operation = 'operation', result?: T) {
    return (error: any): Observable<T> => {

      // TODO: send the error to remote logging infrastructure
      console.error(error); // log to console instead

      // TODO: better job of transforming error for user consumption
      this.log(`${operation} failed: ${error.message}`);

      // Let the app keep running by returning an empty result.
      return of(result as T);
    };
  }
  
  /** Log a HeroService message with the MessageService */
  private log(message: string) {
    //this.messageService.add(`UserService: ${message}`);
  }
}
