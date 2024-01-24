import { Component, OnInit } from '@angular/core';
import {ActivatedRoute} from '@angular/router'
import {DatabaseProjectServices} from "../services/database-project.services";
import {Chart} from "angular-highcharts"
import {openTaskChartOptions} from "./helpers/openTaskChartOptions";
import {ConfigurationService} from "../services/configuration.services";


@Component({
  selector: 'app-compare-project-page',
  templateUrl: './compare-project-page.component.html',
  styleUrls: ['./compare-project-page.component.scss']
})

export class CompareProjectPageComponent implements OnInit {
  projects: string[] = []
  ids: string[] = []
  resultReq: ReqData[] = []
  openTaskChart = new Chart()

  webUrl = ""

  constructor(private configurationService: ConfigurationService,
              private route: ActivatedRoute,
              private dbProjectService: DatabaseProjectServices) {
    this.projects = this.route.snapshot.queryParamMap.getAll("keys")
    this.ids = this.route.snapshot.queryParamMap.getAll("value")
    this.webUrl = configurationService.getValue("webUrl")
  }



  ngOnInit(): void {
    for (let i = 0; i < this.projects.length; i++) {
      this.dbProjectService.getProjectStatByID(this.ids[i]).subscribe(projects => {
        this.resultReq[i] = projects.data
      })
    }

    let colors = ["blue", "green", "red", "orange", "purple", "black"]

    let openTaskElem = document.getElementById('open-task') as HTMLElement;
    let openTaskTitle = document.getElementById('open-task-title') as HTMLElement;
    this.dbProjectService.getComplitedGraph("1", this.projects).subscribe(info => {
      if (info.data["count"] == null) {
        openTaskElem.remove()
        openTaskTitle.remove()
      }
      else{
        // @ts-ignore
        openTaskChartOptions.xAxis["categories"] = info.data["categories"]
        for (let j = 0; j < this.projects.length; j++){
          var count = []
          for (let i = 0; i < info.data["categories"].length; i++){
            // @ts-ignore
            count.push(info.data["count"][info.data["categories"][i]][j])
          }
          openTaskChartOptions.series?.push({ name: this.projects[j],
            type: "column",
            color: colors[j],
            data: count})
          this.openTaskChart = new Chart(openTaskChartOptions)
        }
      }
    })
  }

  ngOnDestroy(): void{
    // @ts-ignore
    openTaskChartOptions.xAxis["categories"] = []
    openTaskChartOptions.series = []
  }
}



class ReqData {
  Id: number;
  Key: string;
  Name: string;
  allIssuesCount: number;
  averageIssuesCount: string;
  averageTime: number;
  closeIssuesCount: number;
  openIssuesCount: number;
  resolvedIssuesCount: number;
  reopenedIssuesCount: number;
  progressIssuesCount: number;
}
