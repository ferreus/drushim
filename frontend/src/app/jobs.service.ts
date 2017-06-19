import { Injectable } from '@angular/core';
import { Http, Response } from '@angular/http';
import { Observable } from 'rxjs/Observable';
import 'rxjs/add/operator/catch';
import 'rxjs/add/operator/map';
import 'rxjs/add/observable/throw';

import { JobItem} from './models';


@Injectable()
export class JobsService {

  constructor(private http: Http) { }

  getJobs() : Observable<JobItem[]> {
    return this.http.get('/v1/jobs').map((res) => {
      return res.json() || [];
    },this.handleError);
  }

  handleError(error: Response | any) : Observable<any> {
    let errMsg: string;
    if (error instanceof Response) {
      let body = <any>{};
      try {
        body = error.json() || '';
      } catch (e) {
        body = {};
      }
      const err = body.error || '';
      errMsg = error.statusText ? `${error.status} - ${error.statusText || ''} ${err}` : 'Backend server error';
    } else {
      errMsg = error.message ? error.message : error.toString();
    }
    console.error(errMsg);
    return Observable.throw(errMsg);
  }

}
