import { Component, OnInit } from '@angular/core';
import { FetchBlockData } from '@lthn/blockchain/letheanservice';

@Component({
  selector: 'app-blockchain',
  standalone: true,
  imports: [],
  template: `<p>Hello from the Blockchain Route!</p>`
})
export class BlockchainComponent implements OnInit {
  ngOnInit() {
    this.fetchData();
  }

  async fetchData() {
    await FetchBlockData();
  }
}
