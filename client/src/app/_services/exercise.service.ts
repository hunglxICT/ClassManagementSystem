import { Injectable } from '@angular/core';
import { HttpClient, HttpHeaders } from '@angular/common/http';

import { environment } from '@environments/environment';
import { Exercise, Submission } from '@app/_models';
import { catchError, map, tap } from 'rxjs/operators';
import { Observable, of } from 'rxjs';

@Injectable({ providedIn: 'root' })
export class ExerciseService {
    
    constructor(private http: HttpClient) { }
    httpOptions = {
    	headers: new HttpHeaders({ 'Content-Type': 'application/json' })
    };
    
    getByID(id: number) {
    	const url = `${environment.apiUrl}/exercise-detail/${id}`;
    	return this.http.get<Exercise>(url).pipe(
      		tap(_ => this.log(`fetched exercise id=${id}`)),
      		catchError(this.handleError<Exercise>(`getByID id=${id}`))
    	);
    }
    
    saveSubmission(submission: Submission) {
        const url = `${environment.apiUrl}/submit`;
        return this.http.post(url, submission, this.httpOptions).pipe(
    		tap(_ => this.log(`add new submission to exercise id=${submission.Exerciseid}`)),
    		catchError(this.handleError<any>('saveSubmission'))
    	);
    }
    
    addLinkSubmission(id: number, submission: any) {
        const url = `${environment.apiUrl}/add-link-submission/${id}`;
        return this.http.post(url, submission).pipe(
    		tap(_ => this.log(`add new link submission to exercise id=${submission.Exerciseid}`)),
    		catchError(this.handleError<any>('addLinkSubmission'))
    	);
    }
    
    
    getSubmissions(id: number) {
        return this.http.get(`${environment.apiUrl}/get-submissions/${id}`).pipe(map((res: any) => {
            return <Submission[]>res['result'];
        }))
    }
    /*
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
    */
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
