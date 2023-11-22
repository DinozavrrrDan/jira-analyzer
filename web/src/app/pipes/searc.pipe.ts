import {Pipe, PipeTransform} from "@angular/core";
import {IProj} from "../models/proj.model";

@Pipe({
  name: 'search_name'
})
export class SearchPipe implements PipeTransform {
    transform(projects: IProj[], value: string) {
      return projects.filter(project =>{
        return project.Name.includes(value);
      })
    }
}