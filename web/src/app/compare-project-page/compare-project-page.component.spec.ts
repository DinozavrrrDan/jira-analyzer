import { ComponentFixture, TestBed } from '@angular/core/testing';

import { CompareProjectPageComponent } from './compare-project-page.component';

describe('CompareProjectPageComponent', () => {
  let component: CompareProjectPageComponent;
  let fixture: ComponentFixture<CompareProjectPageComponent>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      declarations: [ CompareProjectPageComponent ]
    })
    .compileComponents();

    fixture = TestBed.createComponent(CompareProjectPageComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
