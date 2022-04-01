import { Injectable } from '@angular/core';
import { SignupService } from './signup.service';
import { WebRequestService } from './web-request.service';

@Injectable({
  providedIn: 'root'
})
export class SubredditsService {

  getSubreddits() {
    // get data from Backend
    // return this.WebReqService.post()
    return [
      {"name": "Science", "description": "This community is a place to share and discuss new scientific research. Read about the latest advances in astronomy, biology, medicine, physics, social science, and more. Find and submit new publications and popular science coverage of current research."},
      {"name": "Subreddit 2", "description": "empty description"},
      {"name": "Subreddit 3", "description": "empty description"},
      {"name": "Subreddit 4", "description": "empty description"}
    ]
  }

  constructor(private WebReqService: WebRequestService) { }
}
