import {APP_INITIALIZER, NgModule} from '@angular/core';
import { BrowserModule } from '@angular/platform-browser';

import { AppComponent } from './app.component';
import { HomePageComponent } from './home-page/home-page.component';
import {RouterModule} from "@angular/router";
import {ProjectComponent} from "./components/project/project.component";
import {HttpClientModule} from "@angular/common/http";
import { ProjectPageComponent } from './project-page/project-page.component';
import { MyProjectPageComponent } from './my-project-page/my-project-page.component';
import { MyProjectComponent } from './components/my-project/my-project.component';
import {FormsModule} from "@angular/forms";
import {NgxPaginationModule} from "ngx-pagination";
import { ComparePageComponent } from './compare-page/compare-page.component';
import {ProjectWithCheckboxComponent} from "./components/checkbox-with-project/checkbox-with-project.component";
import {CompareProjectPageComponent } from './compare-project-page/compare-project-page.component';
import {ChartModule} from "angular-highcharts";
import {CheckboxWithSettingsComponent} from "./components/checkbox-with-settings/checkbox-with-settings.component";
import { ProjectStatPageComponent } from './project-stat-page/project-stat-page.component';
import { ConfigurationService } from "./services/configuration.services";

export function initApp(configurationService: ConfigurationService) {
  return () => configurationService.load().toPromise();
}


const routes = [
  {path: '', component: HomePageComponent},
  {path: 'projects', component: ProjectPageComponent},
  {path: 'compare', component: ComparePageComponent},
  {path: 'myprojects', component: MyProjectPageComponent},
  {path: 'compare-projects', component: CompareProjectPageComponent},
  {path: 'projects-settings', component: CheckboxWithSettingsComponent},
  {path: 'project-stat', component: ProjectStatPageComponent}
]

@NgModule({
  declarations: [
    AppComponent,
    HomePageComponent,
    ProjectComponent,
    ProjectPageComponent,
    MyProjectComponent,
    MyProjectPageComponent,
    ComparePageComponent,
    ProjectWithCheckboxComponent,
    CompareProjectPageComponent,
    CheckboxWithSettingsComponent,
    ProjectStatPageComponent
  ],
  imports: [
    BrowserModule,
    HttpClientModule,
    RouterModule.forRoot(routes),
    FormsModule,
    NgxPaginationModule,
    ChartModule
  ],
  providers: [
    {
      provide: APP_INITIALIZER,
      useFactory: initApp,
      multi: true,
      deps: [ConfigurationService]
    }
  ],
  bootstrap: [AppComponent]
})
export class AppModule { }
