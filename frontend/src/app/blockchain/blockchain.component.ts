import { Component } from '@angular/core';
// import {FetchBlockData} from '@lthn/blockchain/service';
import { HighchartsChartComponent, ChartConstructorType } from 'highcharts-angular';

@Component({
  selector: 'app-blockchain',
  standalone: true,
  imports: [HighchartsChartComponent],
  template: `<highcharts-chart
    [constructorType]="chartConstructor"
    [options]="chartOptions"
    [(update)]="updateFlag"
    [oneToOne]="oneToOneFlag"
    class="chart"
  />`,
  styles: [`.chart { width: 100%; height: 400px; display: block; }`]
})
export class BlockchainComponent  {
  chartOptions: Highcharts.Options = {
    series: [
      {
        data: [1, 2, 3],
        type: 'line',
      },
    ],
  }
  chartConstructor: ChartConstructorType = 'chart'; // Optional, defaults to 'chart'
  updateFlag: boolean = false; // Optional
  oneToOneFlag: boolean = true;
  // async fetchData() {
  //   await FetchBlockData("0");
  // }
}
