import { ComponentFixture, TestBed } from '@angular/core/testing';

import { ProjectPageComponent } from './project-page.component';
import {ProjectServices} from "../services/project.services";
import {HttpClient} from "@angular/common/http";
import {HttpClientTestingModule} from "@angular/common/http/testing";
import {SearchPipe} from "../pipes/search.pipe";

describe('ProjectPageComponent', () => {
  let component: ProjectPageComponent;
  let httpClient: HttpClient;
  let search_name: SearchPipe;
  let fixture: ComponentFixture<ProjectPageComponent>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      declarations: [
        ProjectPageComponent,
        SearchPipe
      ],
      providers: [
        ProjectServices,
      ],
      imports: [
        HttpClientTestingModule
      ]
    })
    .compileComponents();

    fixture = TestBed.createComponent(ProjectPageComponent);
    httpClient = TestBed.get(HttpClient);
    component = fixture.componentInstance;
    search_name = new SearchPipe();
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
