import { Injectable } from '@angular/core';
import { HttpClient, HttpHeaders } from '@angular/common/http';
import { Observable } from 'rxjs';

@Injectable({
  providedIn: 'root'
})
export class ApiService {
  private apiUrl = 'http://localhost:8001';
  private token = ''

  constructor(private http: HttpClient) { }

  getNotes(): Observable<any> {
    const headers = new HttpHeaders({
      'content-type': 'application/json',
      'Authorization': 'Bearer ' + this.token
    })

    return this.http.get<any>(this.apiUrl + '/notes/get', {
      headers,
    });
  }

  /* getUserById(id: number): Observable<User> {
    return this.http.get<User>(`${this.apiUrl}/${id}`);
  }

  createUser(user: User): Observable<User> {
    return this.http.post<User>(this.apiUrl, user);
  }

  updateUser(user: User): Observable<User> {
    return this.http.put<User>(`${this.apiUrl}/${user.id}`, user);
  }

  deleteUser(id: number): Observable<void> {
    return this.http.delete<void>(`${this.apiUrl}/${id}`);
  } */
}
