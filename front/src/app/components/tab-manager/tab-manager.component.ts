import {Component, EventEmitter, input, Input, Output} from '@angular/core';

@Component({
  selector: 'app-tab-manager',
  standalone: false,
  templateUrl: './tab-manager.component.html',
  styleUrl: './tab-manager.component.css'
})
export class TabManagerComponent {
  @Input({required: true })
  tabs!: string[];

  @Output()
  activeTabChange = new EventEmitter<number>();

  @Input()
  activeTab: number = 0;

  setActiveTab(tab: number): void {
    this.activeTab = tab;
    this.activeTabChange.emit(tab)
  }
}
