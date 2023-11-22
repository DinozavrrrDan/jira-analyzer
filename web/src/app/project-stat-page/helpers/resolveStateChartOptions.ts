import { Options } from 'highcharts';

export const resolveStateChartOptions: Options = {
  legend: {
    enabled: true,
  },
  credits: {
    enabled: false,
  },
  title: {
    text: 'Resolve tasks',
  },
  yAxis: {
    visible: true,
    title: {
      text: 'Resolve issue count'
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

