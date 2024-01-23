import {Injectable} from "@angular/core";
import {HttpClient} from "@angular/common/http";
import {Observable} from "rxjs";
import {IRequest} from "../models/request.model";
import {ConfigurationService} from "./configuration.services";

@Injectable({
    providedIn: 'root'
})
export class ProjectServices {
  urlPath = ""

  constructor(private http: HttpClient, private configurationService: ConfigurationService) {
    this.urlPath = configurationService.getValue("pathUrl")
  }


  // @ts-ignore
  getAll(page: number, searchName: String): Observable<IRequest>{
    return this.http.get<IRequest>('http://'+ this.urlPath +'/api/v1/connector/projects?' +
    'limit=10&page='+page + '&search=' + searchName)
  }

  // @ts-ignore
  addProject(key: String): Observable<IRequest>{
<<<<<<< HEAD
    return this.http.post<IRequest>('http://'+ this.urlPath +'/api/v1/connector/updateProject?project='+key, '')
=======
    return this.http.post<IRequest>('http://127.0.0.1:8003/api/v1/connector/updateProject?project='+key, ' ')
>>>>>>> 77eb92a4ac4e65cecd3e27102698a8f5679caf6c
    // TODO Написать запрос на добавление проета в БД. Добавление происходит по ключу проекта
  }

  // @ts-ignore
  deleteProject(key: String): Observable<IRequest> {
    return this.http.delete<IRequest>('http://'+ this.urlPath +'/api/v1/resource/project?project=' + key)

    // TODO Написать запрос на удаление проекта. Удаление происходит по id проекта в БД.
  }
}
