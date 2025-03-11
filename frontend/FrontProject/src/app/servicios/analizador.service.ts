import { Injectable } from '@angular/core';
import { HttpClient } from '@angular/common/http';
import { Observable } from 'rxjs';

@Injectable({
  providedIn: 'root'
})
export class AnalizadorService {
  private backendUrl = 'http://localhost:8000/analyze'; 

  constructor(private http: HttpClient) { }

  analizarEntrada(entrada: string) {
    return this.http.post<any>(this.backendUrl, { entrada });
  }
}
