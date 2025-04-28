import { Component, EventEmitter, Input, OnInit, Output } from '@angular/core';

@Component({
  selector: 'app-table',
  standalone: false,
  templateUrl: './table.component.html',
  styleUrl: './table.component.css',
})
export class TableComponent<T> implements OnInit {
  @Input({ required: true })
  data!: T[];

  @Input({ required: true })
  columns!: string[];

  @Input()
  pageSize: number = 5;

  @Output()
  click = new EventEmitter<{ row: T; column: string }>();

  currentPage: number = 1;
  totalPages: number = 1;
  paginatedData: any[] = [];

  ngOnInit(): void {
    this.updatePagination();
  }

  updatePagination(): void {
    this.totalPages = Math.ceil(this.data.length / this.pageSize);
    const startIndex = (this.currentPage - 1) * this.pageSize;
    const endIndex = startIndex + this.pageSize;
    this.paginatedData = this.data.slice(startIndex, endIndex);
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
      .replace(/^./, (match) => match.toLowerCase());
    if (anyRow[`${field}Title`]) {
      return anyRow[`${field}Title`];
    } else {
      return anyRow[field];
    }
  }
}
