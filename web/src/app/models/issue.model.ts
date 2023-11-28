import {IProj} from "./proj.model";
import {Time} from "@angular/common";

export class Issue {
  constructor(
 public Key: String,
 public Project: IProj,
 public CreatedTime: Time,
 public ClosedTime: Time,
 public UpdatedTime: Time,
 public Summary: String,
 public Description: String,
 public Type: String,
 public Priority: String,
 public Status: String,
 public Creator: String,
 public Assignee: String,
 public TimeSpent: Number
){}
}
