import {Component, EventEmitter, Input, OnInit, Output} from '@angular/core'
import {IProj} from "../../models/proj.model";
import {ProjectServices} from "../../services/project.services";
import {CheckedProject} from "../../models/check-element.model";

@Component({
  selector: 'project-checkbox',
  templateUrl: './checkbox-with-project.component.html',
  styleUrls: ['./checkbox-with-project.component.css']
})

export class ProjectWithCheckboxComponent implements OnInit {
  @Output() onChecked: EventEmitter<any> = new EventEmitter<{}>();
  @Input() project: IProj
  isChecked: Boolean;

  constructor(private projectService: ProjectServices) {}

  ngOnInit(): void {
    this.isChecked = this.project.Existence;
  }

  changed(isChecked: any) {
    //console.log("Child", this.isChecked, this.project.Name)
    this.onChecked.emit(new CheckedProject(this.project.Name, this.isChecked, this.project.Id));
  }
}

