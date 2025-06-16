import { Injectable } from '@angular/core';
import { HttpClient } from '@angular/common/http';
import { Storage } from './storage';

@Injectable({
  providedIn: 'root'
})
export class ApiService {
  private readonly ROOT_URL = 'http://localhost:8080';

  constructor(private http: HttpClient) { }

  get(uri: string, params?: any) {
    return this.http.get(`${this.ROOT_URL}/${uri}`, { params });
  }

  post(uri: string, payload: object) {
    return this.http.post(`${this.ROOT_URL}/${uri}`, payload);
  }

  patch(uri: string, payload: object) {
    return this.http.patch(`${this.ROOT_URL}/${uri}`, payload);
  }

  delete(uri: string) {
    return this.http.delete(`${this.ROOT_URL}/${uri}`);
  }
}