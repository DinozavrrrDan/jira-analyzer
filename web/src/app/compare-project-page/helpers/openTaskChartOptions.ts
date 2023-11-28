import { Options } from 'highcharts';

export const openTaskChartOptions: Options = {
  chart: {
    type: 'column',
  },
  credits: {
    enabled: false,
  },
  title: {
    text: 'Open task time statistic',
  },
  yAxis: {
    visible: true,
    title: {
      text: 'Open issue count'
    }
  },
  legend: {
    enabled: false,
  },
  xAxis: {
    lineColor: '#fff',
    categories: [],
    title: {
      text: 'Time'
    }
  },

  plotOptions: {
    series: {
      borderRadius: 5,
    } as any,
  },

  series: [
    {
      type: 'column',
      color: '#506ef9',
      data: [],
    },
  ],
};
