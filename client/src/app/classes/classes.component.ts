import { Component, OnInit } from '@angular/core';
import { first } from 'rxjs/operators';
import { Location } from '@angular/common';
import { ActivatedRoute } from '@angular/router';

import { Class } from '@app/_models';
import { ClassService } from '@app/_services';

@Component({
  selector: 'app-classes',
  templateUrl: './classes.component.html',
  styleUrls: ['./classes.component.css']
})
export class ClassesComponent implements OnInit {
  loading = false;
  classes: Class[];
  constructor(private classService: ClassService,
      private route: ActivatedRoute,
      private location: Location
  ) { }

  ngOnInit() {
  	this.loading = true;
  	this.classService.getAll().pipe(first()).subscribe(classes => {
            this.loading = false;
            this.classes = classes;
        });
  }
  
    goBack(): void {
      this.location.back();
    }
    
    delete(classs: Class): void {
      this.classService.deleteByID(classs['Id'])
        .subscribe();
      this.classService.getAll().pipe(first()).subscribe(classes => {
            this.classes = classes;
        });
    }
  
}
