import {
  Component,
  EventEmitter,
  Input,
  Output,
  TemplateRef,
} from '@angular/core';

@Component({
  selector: 'app-tab-manager',
  standalone: false,
  templateUrl: './tab-manager.component.html',
  styleUrl: './tab-manager.component.css',
})
export class TabManagerComponent {
  @Input({ required: true })
  tabs!: { title: string; content?: TemplateRef<any> }[];

  @Output()
  activeTabChange = new EventEmitter<number>();

  activeTab: number = 0;

  setActiveTab(tab: number): void {
    this.activeTab = tab;
    this.activeTabChange.emit(tab);
  }
}
