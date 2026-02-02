import { Injectable } from '@angular/core';
import { HttpClient, HttpHeaders } from '@angular/common/http';
import { Observable } from 'rxjs';

@Injectable({
  providedIn: 'root',
})
export class ApiService {
  private apiUrl = 'http://localhost:8001';
  private token =
    'eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9.eyJuYW1lIjoic2FoaWwiLCJleHAiOjE3MzI1NjU2NDQsImlhdCI6MTczMjM5Mjg0NH0.hLEsU0YQRL3i1xKasncb0A15tqEbKiqnGKKxl1bX5JMTvYZAm0BzUPOpnLxEf383ilg3UXpvOgRkeioRcgvP2hsH2kzx0KfZF07xw3aLvjf1XiXFSKHnGxOa-P20aj6ULKeqPdpIf9U5n1hFXAO---Xfx2ynLUjqq46SUteioweXihoWTPhl3YlloIxrtpbfuBNgHevT1pDOeZfO2IJp19f1fmgZTDlPw7-IUoGk-xC2cZMGdAF35bT7X_kfbn8JwkfjLS2MaiUya34IZYHhVDH5gEj0RSw-7BWG_sTPK6oNhfp5-cUeICA394JqLJevKFaZz18zugCSBMTw6h3MXQ';

  constructor(private http: HttpClient) {}

  getNotes(): Observable<any> {
    const headers = new HttpHeaders({
      'content-type': 'application/json',
      Authorization: 'Bearer ' + this.token,
    });

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
