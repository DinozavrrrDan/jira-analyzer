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


  // @ts-ignore
  getAll(): Observable<IRequest> {
    return this.http.get<IRequest>('http://' + this.urlPath + '/api/v1/resource/projects')
    // TODO Написать запрос на получение всех проектов
  }

  // @ts-ignore
  getProjectStatByID(name: string): Observable<IRequestObject> {
    return this.http.get<IRequest>('http://' + this.urlPath + '/api/v1/resource/project?project=' + name)

    // TODO Написать запрос на получение статистики проекта по ID
  }

  // @ts-ignore
  getComplitedGraph(taskNumber: string, projectName: Array<string>): Observable<IRequestObject> {
    this.http.get<IRequest>('http://' + this.urlPath + '/api/v1/graph/compare/' + taskNumber + '?project=' + projectName.toString())
    // TODO Написать запрос на получение сравнения
  }

  // @ts-ignore
  getGraph(taskNumber: string, projectName: string): Observable<IRequestObject> {
    return this.http.get<IRequest>('http://' + this.urlPath + '/api/v1/graph/get/'+ taskNumber +'?project=' + projectName)
    // TODO Написать запрос на получение графа
  }

  // @ts-ignore
  makeGraph(taskNumber: string, projectName: string): Observable<IRequestObject> {
    return this.http.get<IRequest>('http://' + this.urlPath + '/api/v1/graph/make/'+ taskNumber +'?project=' + projectName)

    // TODO Написать запрос на создание графа
  }

  // @ts-ignore
  deleteGraphs(projectName: string): Observable<IRequestObject> {
    return this.http.get<IRequest>('http://' + this.urlPath + '/api/v1/graph/delete?project=' + projectName)
    // TODO Написать запрос на удаление графа
  }

  // @ts-ignore
  isAnalyzed(projectName: string): Observable<IRequestObject> {
    // TODO Написать запрос
  }

  // @ts-ignore
  isEmpty(projectName: string): Observable<IRequestObject> {
<<<<<<< HEAD:frontend/src/app/services/database-project.services.ts
    // TODO Написать запрос
=======
    return this.http.get<IRequestObject>('http://' + this.urlPath + '/api/v1/isEmpty?project=' + projectName)
>>>>>>> 6a15cb1650a9c1e304607e8f9b48d77b20ebf674:web/src/app/services/database-project.services.ts
  }
}
