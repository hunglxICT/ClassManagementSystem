import { Component, OnInit } from '@angular/core';
import { ActivatedRoute } from '@angular/router';
import { Location } from '@angular/common';

import { Class } from '@app/_models';
import { ClassService } from '@app/_services';

@Component({
  selector: 'app-newclass',
  templateUrl: './newclass.component.html',
  //styleUrls: ['./newclass.component.css']
})
export class NewclassComponent implements OnInit {
  
  classs: Class;
  
  constructor(
    private route: ActivatedRoute,
    private classService: ClassService,
    private location: Location
  ) { }

  ngOnInit() {
    this.classs = new Class;
    //alert(JSON.stringify(this.account));
  }
  
  goBack(): void {
    this.location.back();
  }
  
  save(): void {
    this.classService.addnewclass(this.classs)
      .subscribe(() => this.goBack());
  }
}
