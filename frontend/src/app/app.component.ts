import { Component,OnInit } from '@angular/core';

import { JobsService } from './jobs.service';
import { JobItem } from "app/models";

import * as _ from 'lodash';

@Component({
  selector: 'app-root',
  templateUrl: './app.component.html',
  styleUrls: ['./app.component.css']
})
export class AppComponent implements OnInit {
  title = 'app';
  jobs : JobItem[];
  filteredJobs : JobItem[];
  currentJob : JobItem = null;
  currentIndex : number = 0;
  constructor(private jobsService : JobsService) {}


  private search_ : string = "";
  set search(pattern: string) {
    this.search_ = pattern;
    if (this.search_.length == 0) {
      this.filteredJobs = this.jobs;
    } else {
      this.filterJobs();
    }
  }
  get search() {
    return this.search_;
  }

  ngOnInit(): void {
    this.jobsService.getJobs().subscribe((jobs : JobItem[]) => {
      this.jobs = _.filter(jobs,(j: JobItem) => {
        return j.Regions.length <= 4 && this.matches(j.Regions,'חיפה');
      });
      this.filteredJobs = this.jobs;
      this.currentJob = this.jobs[0];
    })
  }

  selectJob(job: JobItem) {
    this.currentJob = job;
    this.currentIndex = -1;
    for (var i=0;i< this.filteredJobs.length;i++) {
      if (this.filteredJobs[i].JobCode === this.currentJob.JobCode) {
        this.currentIndex = i;
      }
    }
  }

  filterJobs() {
    this.filteredJobs = _.filter(this.jobs,(j:JobItem) => {
      let haifa = j.Regions.length <= 4 && this.matches(j.Regions,'חיפה');
      let pattern = this.search_.toLowerCase();
      let search = j.Name.toLowerCase().includes(pattern) || this.matches(j.Requirements,pattern) 
        || this.matches([j.Experience],pattern) || this.matches(j.Skils,pattern);

      return haifa && search;
    });
  }

  matches(items : string[], pattern: string) : boolean{
    if (items == null) {
      return false;
    }
    for (var i=0;i<items.length;i++) {
      if (items[i].toLowerCase().includes(pattern)) {
        return true;
      }
    }
    return false;
  }

  regions() : string {
    return this.currentJob!=null && this.currentJob.Regions!=null ? this.currentJob.Regions.join(',') : '';
  }

  requirements() : string {
    return this.currentJob!=null && this.currentJob.Requirements!=null ? this.currentJob.Requirements.join('\n') : '';
  }

  skils() : string {
    return this.currentJob!=null && this.currentJob.Skils!=null ? this.currentJob.Skils.join(',') : '';
  }

  scroll(event: KeyboardEvent) {
    if (event.keyCode == 40 && this.currentIndex < this.filteredJobs.length-1) {
      event.preventDefault();
      this.currentJob = this.filteredJobs[++this.currentIndex];
      let tbody = document.getElementById('tbody');
      let current = document.getElementById('item'+this.currentJob.JobCode);
      current.scrollIntoView(false);
    } else if (event.keyCode == 38 && this.currentIndex > 0) {
      event.preventDefault();
      this.currentJob = this.filteredJobs[--this.currentIndex];
      let tbody = document.getElementById('tbody');
      let current = document.getElementById('item'+this.currentJob.JobCode);
      current.scrollIntoView(false);
    }
  }


}
