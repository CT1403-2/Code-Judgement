import { Component, EventEmitter, Input, OnInit, Output } from '@angular/core';

@Component({
  selector: 'app-table',
  standalone: false,
  templateUrl: './table.component.html',
  styleUrl: './table.component.css'
})
export class TableComponent<T> implements OnInit {
  @Input({ required: true })
  data!: T[];

  @Input({ required: true })
  columns!: string[];

  @Input({ required: true })
  totalPages!: number;

  @Input()
  actionIcon?: string;

  @Output()
  click = new EventEmitter<{ row: T; column: string }>();

  @Output()
  action = new EventEmitter<T>();

  @Output()
  pageChange = new EventEmitter<number>();

  currentPage: number = 1;

  ngOnInit(): void {
    this.updatePagination();
  }

  updatePagination(): void {
    this.pageChange.emit(this.currentPage);
  }

  nextPage(): void {
    if (this.currentPage < this.totalPages) {
      this.currentPage++;
      this.updatePagination();
    }
  }

  previousPage(): void {
    if (this.currentPage > 1) {
      this.currentPage--;
      this.updatePagination();
    }
  }

  getData(row: T, column: string) {
    const anyRow = row as any;
    const field = column
      .replace(/\s+/g, '')
      .replace(/^./, match => match.toLowerCase());
    if (anyRow[`${field}Title`]) {
      return anyRow[`${field}Title`];
    } else {
      return anyRow[field];
    }
  }
}
