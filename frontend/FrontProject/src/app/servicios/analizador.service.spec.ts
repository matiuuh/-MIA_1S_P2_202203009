import { Injectable } from '@angular/core';
import { HttpClient } from '@angular/common/http';
import { Observable } from 'rxjs';

@Injectable({
  providedIn: 'root'
})
export class AnalizadorService {
  private apiUrl = 'http://localhost:8000/analyze'; // URL del backend

  constructor(private http: HttpClient) {}

  analizarEntrada(entrada: string): Observable<{ resultado: string }> {
    const body = { entrada };  // JSON con la entrada
    return this.http.post<{ resultado: string }>(this.apiUrl, body);
  }
}
