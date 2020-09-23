import { Injectable } from '@angular/core';
import { HttpClient, HttpHeaders } from '@angular/common/http';

import { catchError, map, tap } from 'rxjs/operators';
import { Observable, of } from 'rxjs';

import { environment } from '@environments/environment';
import { Class, Exercise } from '@app/_models';

@Injectable({ providedIn: 'root' })
export class ClassService {
    constructor(private http: HttpClient) { }
    httpOptions = {
    	headers: new HttpHeaders({ 'Content-Type': 'application/json' })
    };
    
    httpOptions2 = {
    	headers: new HttpHeaders({ 'Content-Type': 'undefined' })
    };
    getAll() {
        return this.http.get(`${environment.apiUrl}/list-class`).pipe(map((res: any) => {
            //res['playload'] = res;
            //console.log('~'+JSON.stringify(res['result'])+'~');
            return <Class[]>res['result'];
        }))
    }
    
    getByID(id: number) {
    	const url = `${environment.apiUrl}/class-info/${id}`;
    	return this.http.get<Class>(url).pipe(
      		tap(_ => this.log(`fetched class id=${id}`)),
      		catchError(this.handleError<Class>(`getByID id=${id}`))
    	);
    }
    
    deleteByID(id: number) {
    	const url = `${environment.apiUrl}/delete-class/${id}`;
    	return this.http.get(url).pipe(
      		tap(_ => this.log(`delete account id=${id}`)),
      		catchError(this.handleError(`deleteByID id=${id}`))
    	);
    }
    
    getEnrolledStudents(id: number) {
    	const url = `${environment.apiUrl}/enrolled-students/${id}`;
    	return this.http.get(url).pipe(
      		tap(_ => this.log(`get enrolled for class with id=${id}`)),
      		catchError(this.handleError(`getEnrolledStudents id=${id}`))
    	);
    }
    
    joinclass(classid: number, studentid: number) {
        const url = `${environment.apiUrl}/join-class`;
        return this.http.post(url, {"classid":classid, "studentid":studentid}, this.httpOptions).pipe(
    		tap(_ => this.log(`student id=${studentid} join class id=${classid}`)),
    		catchError(this.handleError<any>('joinclass'))
    	);
    }
    
    dismiss(classid: number, studentid: number) {
        const url = `${environment.apiUrl}/dismiss-student`;
        return this.http.post(url, {"classid":classid, "studentid":studentid}, this.httpOptions).pipe(
    		tap(_ => this.log(`student id=${studentid} has been dismissed from class id=${classid}`)),
    		catchError(this.handleError<any>('dismiss'))
    	);
    }
    
    addExercise(exercise: Exercise) {
        const url = `${environment.apiUrl}/add-exercise`;
        return this.http.post(url, exercise).pipe(
    		tap(_ => this.log(`add new exercise to class id=${exercise.Classid}`)),
    		catchError(this.handleError<any>('addexercise'))
    	);
    }
    
    addLinkExercise(id: number, exercise: any) {
        const url = `${environment.apiUrl}/add-link-exercise/${id}`;
        return this.http.post(url, exercise).pipe(
    		tap(_ => this.log(`add new link exercise to class id=${exercise.Classid}`)),
    		catchError(this.handleError<any>('addexercise'))
    	);
    }
    
    getExercises(id: number) {
    	const url = `${environment.apiUrl}/list-exercises/${id}`;
    	return this.http.get(url).pipe(
      		tap(_ => this.log(`get exercises for class with id=${id}`)),
      		catchError(this.handleError(`getExercises id=${id}`))
    	);
    }
    
    /*saveByID(account: UserBackend): Observable<UserBackend> {
    	const url = `${environment.apiUrl}/edit-profile`;
    	return this.http.post<UserBackend>(url, account, this.httpOptions).pipe(
    		tap(_ => this.log(`updated account id=${account.Id}`)),
    		catchError(this.handleError<any>('saveByID'))
    	);
    }*/
    
    addnewclass(classs: Class): Observable<Class> {
    	const url = `${environment.apiUrl}/new-class`;
    	return this.http.post<Class>(url, classs, this.httpOptions).pipe(
    		tap(_ => this.log(`add new class id=${classs.Id}`)),
    		catchError(this.handleError<any>('addnewclass'))
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
