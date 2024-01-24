import { Component } from '@angular/core';
import {ConfigurationService} from "./services/configuration.services";

@Component({
  selector: 'app-root',
  templateUrl: './app.component.html',
  styleUrls: ['./app.component.css']
})
export class AppComponent {
  title = 'web';
  webUrl = ""
  constructor(private configurationService: ConfigurationService) {
    this.webUrl = configurationService.getValue("webUrl")
  }

}
