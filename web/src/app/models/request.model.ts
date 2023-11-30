import {Links} from "./links.model";
import {PageInfo} from "./pageInfo.model";

export class IRequest {
  constructor(
  public _links: Links,
  public data: [],
  public message: String,
  public name: String,
  public pageInfo: PageInfo,
  public status: Boolean){}
}


