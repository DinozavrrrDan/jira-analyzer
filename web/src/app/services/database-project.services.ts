import {Injectable} from "@angular/core";
import {HttpClient} from "@angular/common/http";
import {Observable} from "rxjs";
import {IRequest} from "../models/request.model";
import {IRequestObject} from "../models/requestObj.model";
import {ConfigurationService} from "./configuration.services";

@Injectable({
  providedIn: 'root'
})
export class DatabaseProjectServices {
  urlPath = ""

  constructor(private http: HttpClient, private configurationService: ConfigurationService) {
    this.urlPath = configurationService.getValue("pathUrl")
  }


  getAll(): Observable<IRequest>{
    return this.http.get<IRequest>('http://'+ this.urlPath + '/api/v1/projects')
  }

  getProjectStatByID(id: string): Observable<IRequestObject> {
    return this.http.get<IRequestObject>('http://'+ this.urlPath + '/api/v1/projects/'+id)
  }

  getComplitedGraph(taskNumber: string, projectName: Array<string>): Observable<IRequestObject> {
    let projectsString = 'http://'+ this.urlPath + '/api/v1/compare/'+taskNumber+'?project='+projectName.toString()
    return this.http.get<IRequestObject>(projectsString)
  }

  getGraph(taskNumber: string, projectName: string): Observable<IRequestObject> {
    let projectsString = 'http://'+ this.urlPath + '/api/v1/graph/get/'+taskNumber+'?project='+projectName
    return this.http.get<IRequestObject>(projectsString)
  }

  makeGraph(taskNumber: string, projectName: string): Observable<IRequestObject> {
    let projectsString = 'http://'+ this.urlPath + '/api/v1/graph/make/'+taskNumber+'?project='+projectName
    return this.http.post<IRequestObject>(projectsString, ' ')
  }

  deleteGraphs(projectName: string): Observable<IRequestObject> {
    let projectsString = 'http://'+ this.urlPath + '/api/v1/graph/delete'+'?project='+projectName
    return this.http.delete<IRequestObject>(projectsString)
  }

  isAnalyzed(projectName: string): Observable<IRequestObject>{
    let projectsString = 'http://'+ this.urlPath + '/api/v1/isAnalyzed?project='+projectName
    return this.http.get<IRequestObject>(projectsString)
  }

  isEmpty(projectName: string): Observable<IRequestObject>{
    let projectsString = 'http://'+ this.urlPath + '/api/v1/isEmpty?project='+projectName
    return this.http.get<IRequestObject>(projectsString)
  }
}