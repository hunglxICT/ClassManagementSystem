import { Component, OnInit } from '@angular/core';
import { ActivatedRoute } from '@angular/router';
import { Location } from '@angular/common';
import { first } from 'rxjs/operators';
import { FormGroup, FormControl, Validators, FormBuilder } from '@angular/forms';

import { Exercise, Submission } from '@app/_models';
import { ExerciseService } from '@app/_services';

@Component({
  selector: 'app-exercise-detail',
  templateUrl: './exercise-detail.component.html',
  styleUrls: ['./exercise-detail.component.css']
})
export class ExerciseDetailComponent implements OnInit {

  exercise: Exercise;
  newsubmission: Submission;
  uploadForm: FormGroup;
  submissions: Submission[];

  constructor(
    private route: ActivatedRoute,
    private exerciseService: ExerciseService,
    private location: Location,
    private formBuilder: FormBuilder
  ) {}

  ngOnInit(): void {
    this.uploadForm = this.formBuilder.group({
      profile: ['']
    });
    this.getExercise();
  }
  
  onFileSelect(event) {
    if (event.target.files.length > 0) {
      const file = event.target.files[0];
      this.uploadForm.get('profile').setValue(file);
    }
  }
  
  getExercise(): void {
    const id = +this.route.snapshot.paramMap.get('id');
    this.exerciseService.getByID(id)
      .subscribe(exer => {
      this.exercise = exer['result'];
      this.getSubmissions(id);
      });
  }

  goBack(): void {
    this.location.back();
  }
  
  initsubmit(): void {
    this.newsubmission = new Submission;
  }
  
  getSubmissions(classid: number): void {
    this.exerciseService.getSubmissions(classid).subscribe(result => {
        this.submissions = result;
    })
  }
  
  saveSubmission(exercise: Exercise, submission: Submission): void {
    const exerciseid = exercise['Id'];
    submission.Exerciseid = exerciseid;
    const formData = new FormData();
    formData.append('file', this.uploadForm.get('profile').value);
    var id = -1;
    this.exerciseService.saveSubmission(submission).subscribe(result => {
        id = result['result'];
        this.exerciseService.addLinkSubmission(id, formData).subscribe();
    })
  }
  
}
