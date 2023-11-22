import { ComponentFixture, TestBed } from '@angular/core/testing';

import { MyProjectPageComponent } from './my-project-page.component';
import {HttpClient} from "@angular/common/http";
import {ProjectServices} from "../services/project.services";
import {HttpClientTestingModule} from "@angular/common/http/testing";

describe('MyProjectPageComponent', () => {
  let component: MyProjectPageComponent;
  let httpClient: HttpClient;
  let fixture: ComponentFixture<MyProjectPageComponent>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      declarations: [
        MyProjectPageComponent
      ],
      providers: [
        MyProjectPageComponent
      ],
      imports: [
        HttpClientTestingModule
      ]
    })
    .compileComponents();

    fixture = TestBed.createComponent(MyProjectPageComponent);
    httpClient = TestBed.get(HttpClient);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});