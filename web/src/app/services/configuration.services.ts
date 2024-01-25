import { Injectable } from '@angular/core';
import {HttpClient} from "@angular/common/http";
import {mapTo, Observable, tap} from "rxjs";

@Injectable({
  providedIn: 'root'
})
export class ConfigurationService {

  private configuration = {};

  constructor(private httpClient: HttpClient) {
  }

  load(): Observable<void> {
    return this.httpClient.get('/assets/config.json')
      .pipe(
        tap((configuration: any) => this.configuration = configuration),
        mapTo(undefined),
      );
  }

  getValue(key: string, defaultValue?: any): any {
    // @ts-ignore
    return this.configuration[key] || defaultValue;
  }

}
