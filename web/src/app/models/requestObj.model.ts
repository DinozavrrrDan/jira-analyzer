import {Links} from "./links.model";
import {PageInfo} from "./pageInfo.model";

export class IRequestObject {
  constructor(
    public _links: Links,
    public data: any,
    public message: String,
    public name: String,
    public pageInfo: PageInfo,
    public status: Boolean){}
}


