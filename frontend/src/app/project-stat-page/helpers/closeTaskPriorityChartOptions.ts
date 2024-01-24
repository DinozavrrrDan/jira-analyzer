import { Options } from 'highcharts';

export const closeTaskPriorityChartOptions: Options = {
  chart: {
    type: 'column',
  },
  credits: {
    enabled: false,
  },
  title: {
    text: 'CLose task priority',
  },
  yAxis: {
    visible: true,
    title: {
      text: 'Close issue count'
    }
  },
  legend: {
    enabled: false,
  },
  xAxis: {
    lineColor: '#fff',
    categories: [],
    title: {
      text: 'Priority'
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
