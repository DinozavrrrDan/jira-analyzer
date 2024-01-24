import { Options } from 'highcharts';

export const reopenStateChartOptions: Options = {
  legend: {
    enabled: true,
  },
  credits: {
    enabled: false,
  },
  title: {
    text: 'Reopen tasks',
  },
  yAxis: {
    visible: true,
    title: {
      text: 'Reopen issue count'
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

