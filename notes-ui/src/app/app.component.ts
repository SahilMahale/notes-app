import { Component, OnInit } from '@angular/core';
import { RouterOutlet } from '@angular/router';
import { CardComponent } from './card/card.component';
import { ApiService } from './api/api.service';

@Component({
  selector: 'app-root',
  standalone: true,
  imports: [RouterOutlet, CardComponent],
  templateUrl: './app.component.html',
  styleUrl: './app.component.css'
})
export class AppComponent implements OnInit {
  title = 'notes-ui';
  noteTitle = "Sahil's First note"
  noteBody = "Note Body is ochinchin"
  notes: any
  constructor(private apiService: ApiService) { }
  ngOnInit() {
    this.fetchNotes()
  }
  fetchNotes() {
    this.apiService.getNotes().subscribe(
      {
        next: (notes) => {
          this.notes = notes
        }
      }
    )
  }
}
