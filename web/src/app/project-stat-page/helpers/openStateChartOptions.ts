import { Options } from 'highcharts';

export const openStateChartOptions: Options = {
  legend: {
    enabled: true,
  },
  credits: {
    enabled: false,
  },
  title: {
    text: 'Open tasks',
  },
  yAxis: {
    visible: true,
    title: {
      text: 'Open issue count'
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

