import { Component } from '@angular/core';
import { RouterOutlet } from '@angular/router';
import { FormsModule } from '@angular/forms'; // Importa FormsModule


@Component({
  selector: 'app-root',
  imports: [FormsModule], // Agrega FormsModule a los imports
  templateUrl: './app.component.html',
  styleUrl: './app.component.css'
})
export class AppComponent {
  title = 'FrontProject';
  entrada: string = '';
  consola: string = '';

  ejecutar() {
    this.consola = `Resultado: ${this.entrada}`;
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
