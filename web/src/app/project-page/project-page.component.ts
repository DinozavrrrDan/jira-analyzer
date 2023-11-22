import { Component, OnInit } from '@angular/core';
import {ProjectServices} from "../services/project.services";
import {IProj} from "../models/proj.model";
import {PageInfo} from "../models/pageInfo.model";

@Component({
  selector: 'app-project-page',
  templateUrl: './project-page.component.html',
  styleUrls: ['./project-page.component.css']
})
export class ProjectPageComponent implements OnInit {
  projects: IProj[] = []
  loading = false
  searchName = ''
  pageInfo: PageInfo
  start_page = 1

  constructor(private projectService: ProjectServices) {
  }

  ngOnInit(): void {
    this.loading = true
    this.projectService.getAll(this.start_page, this.searchName).subscribe(projects => {
      this.projects = projects.data
      this.loading = false
      this.pageInfo = projects.pageInfo
    },
      error => {
        // TODO Обработать ошибки от сервера
      })
  }

  gty(page: any){
    this.projectService.getAll(page, this.searchName).subscribe(projects => {
      this.projects = projects.data
      this.pageInfo = projects.pageInfo
      this.loading = false
    })
  }

  getSearchProjects() {
    this.pageInfo.currentPage = this.start_page;
    this.gty(this.pageInfo.currentPage);
  }
}
