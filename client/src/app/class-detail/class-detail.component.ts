import { Component, OnInit } from '@angular/core';
import { ActivatedRoute } from '@angular/router';
import { Location } from '@angular/common';
import { first } from 'rxjs/operators';

import { FormGroup, FormControl, Validators, FormBuilder } from '@angular/forms';

import { User, UserBackend, Class, Exercise } from '@app/_models';
import { UserService, ClassService } from '@app/_services';
import { FileUploadComponent } from '../app-file-upload/app-file-upload.component';

@Component({
  selector: 'app-class-detail',
  templateUrl: './class-detail.component.html',
  styleUrls: ['./class-detail.component.css']
})
export class ClassDetailComponent implements OnInit {
  
  classs: Class;
  teacher: UserBackend;
  students: UserBackend[];
  users: UserBackend[];
  exercises: Exercise[];
  newexercises: Exercise;
  uploadForm: FormGroup;
  
  constructor(
    private route: ActivatedRoute,
    private classService: ClassService,
    private accountService: UserService,
    private location: Location,
    private formBuilder: FormBuilder
  ) { }

  ngOnInit(): void {
    this.uploadForm = this.formBuilder.group({
      profile: ['']
    });
    this.getClass();
  }
  
  goBack(): void {
    this.location.back();
  }
  
  onFileSelect(event) {
    if (event.target.files.length > 0) {
      const file = event.target.files[0];
      this.uploadForm.get('profile').setValue(file);
    }
  }
  
  getClass(): void {
    const id = +this.route.snapshot.paramMap.get('id');
    this.classService.getByID(id)
      .subscribe(cla => {
        this.classs = cla['result'];
        this.getTeacherByClassID(this.classs.Teacherid);
        this.getEnrolledStudents(this.classs.Id);
        this.getExercises(this.classs.Id);
      //alert(JSON.stringify(cla));
      });
  }
  
  getTeacherByClassID(teacherid: number): void {
    this.accountService.getByID(teacherid)
      .subscribe(tea => this.teacher = tea['result']);
  }
  
  getEnrolledStudents(classid: number): void {
    this.classService.getEnrolledStudents(classid)
      .subscribe(stus => this.students = stus['result']);
  }
  
  getExercises(classid: number): void {
    this.classService.getExercises(classid)
      .subscribe(exer => this.exercises = exer['result']);
  }
  
  listtojoinclass(): void {
    this.accountService.getAll().subscribe(users => this.users = users);
  }
  
  closelisttojoinclass(): void {
    if (this.users) delete this.users;
  }
  
  joinclass(classs: Class, account: UserBackend): void {
    const classid = classs['Id'];
    const accountid = account['Id'];
    this.classService.joinclass(classid, accountid).subscribe();
    this.getEnrolledStudents(this.classs.Id);
  }
  
  dismiss(classs: Class, account: UserBackend): void {
    const classid = classs['Id'];
    const accountid = account['Id'];
    this.classService.dismiss(classid, accountid).subscribe();
    this.getEnrolledStudents(this.classs.Id);
  }
  
  initaddexercise(): void {
    this.newexercises = new Exercise;
  }
  
  closeaddexercise(): void {
    if (this.newexercises) delete this.newexercises;
  }
  
  saveexercise(classs: Class, exercise: Exercise): void {
    const classid = classs['Id'];
    exercise.Classid = classid;
    const formData = new FormData();
    formData.append('file', this.uploadForm.get('profile').value);
    //formData.append('data', JSON.stringify(exercise));
    //const uploadfile = this.uploadForm.get('profile').value;
    var id = -1;
    this.classService.addExercise(exercise).subscribe(result => {
        id = result['result'];
        this.classService.addLinkExercise(id, formData).subscribe(_ => this.getExercises(this.classs.Id));
    });
    
  }
  
  deleteexercise(): void {
  }

}

