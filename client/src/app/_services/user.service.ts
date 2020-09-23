import { Injectable } from '@angular/core';
import { HttpClient, HttpHeaders } from '@angular/common/http';

import { environment } from '@environments/environment';
import { User, UserBackend } from '@app/_models';
import { catchError, map, tap } from 'rxjs/operators';
import { Observable, of } from 'rxjs';

@Injectable({ providedIn: 'root' })
export class UserService {
    constructor(private http: HttpClient) { }
    httpOptions = {
    	headers: new HttpHeaders({ 'Content-Type': 'application/json' })
    };
    getAll() {
        return this.http.get(`${environment.apiUrl}/list-student`).pipe(map((res: any) => {
            //res['playload'] = res;
            //console.log('~'+JSON.stringify(res['result'])+'~');
            return <UserBackend[]>res['result'];
        }))
    }
    
    getByID(id: number) {
    	const url = `${environment.apiUrl}/profile/${id}`;
    	return this.http.get<UserBackend>(url).pipe(
      		tap(_ => this.log(`fetched account id=${id}`)),
      		catchError(this.handleError<UserBackend>(`getByID id=${id}`))
    	);
    }
    
    deleteByID(id: number) {
    	const url = `${environment.apiUrl}/delete-account/${id}`;
    	return this.http.get(url).pipe(
      		tap(_ => this.log(`delete account id=${id}`)),
      		catchError(this.handleError(`deleteByID id=${id}`))
    	);
    }
    
    saveByID(account: UserBackend): Observable<UserBackend> {
    	const url = `${environment.apiUrl}/edit-profile`;
    	return this.http.post<UserBackend>(url, account, this.httpOptions).pipe(
    		tap(_ => this.log(`updated account id=${account.Id}`)),
    		catchError(this.handleError<any>('saveByID'))
    	);
    }
    
    register(account: UserBackend): Observable<UserBackend> {
    	const url = `${environment.apiUrl}/new-student`;
    	return this.http.post<UserBackend>(url, account, this.httpOptions).pipe(
    		tap(_ => this.log(`add new account id=${account.Id}`)),
    		catchError(this.handleError<any>('register'))
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
