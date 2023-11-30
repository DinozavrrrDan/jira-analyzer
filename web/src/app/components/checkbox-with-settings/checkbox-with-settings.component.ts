import {Component, EventEmitter, Input, OnInit, Output} from '@angular/core'
import {ProjectServices} from "../../services/project.services";
import {CheckedSetting} from "../../models/check-setting.model";
import {SettingBox} from "../../models/setting.model";

@Component({
  selector: 'projects-settings',
  templateUrl: './checkbox-with-settings.component.html',
  styleUrls: ['./checkbox-with-settings.component.css']
})

export class CheckboxWithSettingsComponent implements OnInit {
  @Output() onChecked: EventEmitter<any> = new EventEmitter<{}>();
  @Input() box: SettingBox
  isChecked: Boolean;

  constructor(private projectService: ProjectServices) {}

  ngOnInit(): void {
    this.isChecked = this.box.Checked;
  }

  changed(isChecked: any) {
    console.log("Child", this.box.Text, this.isChecked, this.box.BoxId)
    this.onChecked.emit(new CheckedSetting(this.box.Text, this.isChecked, this.box.BoxId));
  }
}

