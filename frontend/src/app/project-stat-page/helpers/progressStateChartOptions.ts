import { Options } from 'highcharts';

export const progressStateChartOptions: Options = {
  legend: {
    enabled: true,
  },
  credits: {
    enabled: false,
  },
  title: {
    text: 'Progress tasks',
  },
  yAxis: {
    visible: true,
    title: {
      text: 'Progress issue count'
    }
  },

  xAxis: {
    visible: true,
    categories: [],
    title: {
      text: 'Time'
    }
  },


  series: [],
};

