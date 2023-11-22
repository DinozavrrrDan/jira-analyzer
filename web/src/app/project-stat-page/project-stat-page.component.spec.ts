import { ComponentFixture, TestBed } from '@angular/core/testing';

import { ProjectStatPageComponent } from './project-stat-page.component';

describe('ProjectStatPageComponent', () => {
  let component: ProjectStatPageComponent;
  let fixture: ComponentFixture<ProjectStatPageComponent>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      declarations: [ ProjectStatPageComponent ]
    })
    .compileComponents();

    fixture = TestBed.createComponent(ProjectStatPageComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
