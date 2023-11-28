import {Component, EventEmitter, Input, OnInit, Output} from '@angular/core';
import {IProj} from "../../models/proj.model";
import {CheckedProject} from "../../models/check-element.model";
import {CheckedSetting} from "../../models/check-setting.model";
import {SettingBox} from "../../models/setting.model";
import {Router} from "@angular/router";
import {DatabaseProjectServices} from "../../services/database-project.services";
import {delay} from "rxjs";

@Component({
  selector: 'app-my-project',
  templateUrl: './my-project.component.html',
  styleUrls: ['./my-project.component.css']
})
export class MyProjectComponent implements OnInit{
  @Output() onChecked: EventEmitter<any> = new EventEmitter<{}>();
  @Input() myProject: IProj
  stat: ProjectStat = new ProjectStat()
  checked = 0
  complited = 0
  processed: boolean
  settings: boolean
  checkboxes: SettingBox[] = []
  setting: Map<any, any> = new Map();

  constructor(private router: Router, private dbProjectService: DatabaseProjectServices) {
  }

  ngOnInit(): void{
    this.processed=false;
    this.settings = false;
    this.checkboxes.push(new SettingBox("Гистограмма, отражающая время, которое задачи провели в открытом состоянии", false, 1 ))
    this.checkboxes.push(new SettingBox("Диаграммы, которые показывают распределение времени по состоянием задач", false, 2 ))
    this.checkboxes.push(new SettingBox("График активности по задачам", false, 3 ))
    this.checkboxes.push(new SettingBox("График сложности задач", false, 4 ))
    this.checkboxes.push(new SettingBox("График, отражающий приоритетность всех задач", false, 5 ))
    this.checkboxes.push(new SettingBox("График, отражающий приоритетность закрытых задач", false, 6 ))

    this.dbProjectService.getProjectStatByID(this.myProject.Id.toString()).subscribe(projects => {

      this.stat.AverageIssuesCount = projects.data["allIssuesCount"]
      this.stat.OpenIssuesCount = projects.data["openIssuesCount"]
      this.stat.AllIssuesCount = projects.data["allIssuesCount"]
      this.stat.AverageTime = projects.data["averageTime"]
      this.stat.CloseIssuesCount = projects.data["closeIssuesCount"]
      this.stat.ReopenedIssuesCount = projects.data["reopenedIssuesCount"]
      this.stat.ResolvedIssuesCount = projects.data["resolvedIssuesCount"]
      this.stat.ProgressIssuesCount = projects.data["progressIssuesCount"]

      console.log(projects.data)
      console.log(this.stat)
    })

    this.dbProjectService.isAnalyzed(this.myProject.Name.toString()).subscribe(info => {
      if (info.data["isAnalyzed"]){
        this.processed = true
      }
    })
  }

  async processProject() {
    this.dbProjectService.isEmpty(this.myProject.Name.toString()).subscribe(resp => {
      if (resp.data["isEmpty"]){
        alert("Project is empty, analytical tasks are not available")
        return
      }
      this.dbProjectService.deleteGraphs(this.myProject.Name.toString()).subscribe(info => {
        this.complited = 0;
        this.checked = 0;
        this.checkboxes.forEach((box: SettingBox) => {
          if (box.Checked) {
            this.checked++
          }
        })
        for (const box of this.checkboxes) {
          if (box.Checked) {
            this.dbProjectService.makeGraph(box.BoxId.toString(), this.myProject.Name.toString()).subscribe(info => {
              this.complited++;
            })
          }
        }
        this.processed = true
      })
    })
  }

  checkResult(){
    let ids:  number[] = []
    let items = this.myProject.Name

    this.checkboxes.forEach((box: SettingBox) =>{
      if (box.Checked){
        ids.push(Number(box.BoxId))
      }
    })

    this.router.navigate([`/project-stat`], {
      queryParams: {
        keys: items,
        value: ids
      }
    });
  }

  clickOnSettings(){
    this.settings = !this.settings;
  }

  disableCheckResultButton(){
    return !this.processed || this.checked != this.complited;
  }

  disableAnalyzeButton(){
    return !this.checkboxes.some(checkbox => checkbox.Checked)  || this.checked != this.complited;
  }

  childOnChecked(setting: CheckedSetting){
    if (setting.Checked) {
      this.setting.set(setting.ProjectName, setting.BoxId)
    }else if (this.setting.has(setting.ProjectName)){
      this.setting.delete(setting.ProjectName)
    }
    this.checkboxes[Number(setting.BoxId) - 1].Checked = setting.Checked
    console.log("Parent ", setting.ProjectName, setting.BoxId, this.checkboxes[Number(setting.BoxId) - 1].Checked)
  }

}

class ProjectStat {
  AllIssuesCount: number;
  AverageIssuesCount: number;
  AverageTime: number;
  CloseIssuesCount: number;
  OpenIssuesCount: number;
  ResolvedIssuesCount: number;
  ReopenedIssuesCount: number;
  ProgressIssuesCount: number
}
