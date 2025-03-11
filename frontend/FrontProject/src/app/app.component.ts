import { Component } from '@angular/core';
import { FormsModule } from '@angular/forms';
import { HttpClientModule } from '@angular/common/http';
import { AnalizadorService } from './servicios/analizador.service';

@Component({
  selector: 'app-root',
  standalone: true,
  imports: [FormsModule, HttpClientModule],
  templateUrl: './app.component.html',
  styleUrls: ['./app.component.css']
})
export class AppComponent {
  title = 'FrontProject';
  entrada: string = '';
  consola: string = '';

  constructor(private analizadorService: AnalizadorService) {}

  ejecutar() {
    this.analizadorService.analizarEntrada(this.entrada).subscribe({
      next: (resp: any) => {
        this.consola = resp.resultado;
      },
      error: (err) => {
        this.consola = `Error: ${err.message}`;
      }
    });
  }

  limpiar() {
    this.entrada = '';
    this.consola = '';
  }

  cargarArchivo(event: any) {
    const file: File = event.target.files[0];
    if (file) {
      const reader = new FileReader();
      reader.onload = () => {
        this.entrada = reader.result as string;
      };
      reader.readAsText(file);
    }
  }
}