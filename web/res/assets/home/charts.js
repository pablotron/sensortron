(() => {
  'use strict';

  // chart options
  const OPTIONS = [{
    id: 't',
    options: {
      maintainAspectRatio: false,

      scales: {
        y: {
          title: {
            display: true,
            text: 'Temperature (°F)',
          },

          grid: { display: false },
        },

        x: {
          type: 'time',

          title: {
            display: true,
            text: 'Time',
          },

          ticks: {
            // limit maximum number of ticks
            maxTicksLimit: 8,
          },

          grid: { display: false },
          time: {
            minUnit: 'hour',

            // ref: https://date-fns.org/v3.6.0/docs/format
            tooltipFormat: 'MM/dd HH:mm',
          },
        },
      },

      plugins: {
        title: {
          display: true,
          text: 'Temperature',
          font: { weight: 'bold', size: 18 },
        },

        tooltip: {
          mode: 'index',
          intersect: false,
        },
      },
    },
  }, {
    id: 'h',
    options: {
      maintainAspectRatio: false,

      scales: {
        y: {
          title: {
            display: true,
            text: 'Humidity (%)',
          },

          grid: { display: false },
        },

        x: {
          type: 'time',

          title: {
            display: true,
            text: 'Time',
          },

          ticks: {
            // limit maximum number of ticks
            maxTicksLimit: 8,
          },

          grid: { display: false },
          time: {
            minUnit: 'hour',

            // ref: https://date-fns.org/v3.6.0/docs/format
            tooltipFormat: 'MM/dd HH:mm',
          },
        },
      },

      plugins: {
        title: {
          display: true,
          text: 'Humidity',
          font: { weight: 'bold', size: 18 },
        },

        tooltip: {
          mode: 'index',
          intersect: false,
        },
      },
    },
  }];

  // download file with given name and data
  const download = (name, data) => {
    let a = document.createElement('a');
    a.download = name;
    a.href = data;
    document.body.appendChild(a);
    a.click();
    document.body.removeChild(a);
    a = null;
  };

  // poll for current chart data
  const poll = () => fetch('/api/home/charts/poll', { method: 'POST' }).then(
    (r) => r.json()
  ).then((r) => {
    // get selected time filter, convert to slice start
    const start = 4 * -document.querySelector('input.time[type="radio"]:checked').value;

    for (let k of Object.keys(r)) {
      if (k in charts) {
        // filter labels and datasets
        r[k].labels = r[k].labels.slice(start);
        r[k].datasets = r[k].datasets.map((set) => {
          set.data = set.data.slice(start);
          return set;
        });

        // update chart
        charts[k].data = r[k];
        charts[k].update();
      }
    }
  });

  // time filter btn click event handler
  document.querySelectorAll('input.time, label.time').forEach((e) => {
    e.addEventListener('click', () => { setTimeout(poll, 10) })
  });

  // chart download btn click event handler
  document.querySelectorAll('.chart-download').forEach((e) => {
    e.addEventListener('click', () => {
      download(e.dataset.name, charts[e.dataset.id].toBase64Image());
      setTimeout(() => document.body.click(), 10);
    });
  });

  // bind to saved event
  document.getElementById('edit-dialog').addEventListener('saved', poll);

  // init charts
  Chart.defaults.color = '#eee';
  const charts = OPTIONS.reduce((r, {id, options}) => {
    r[id] = new Chart(document.getElementById('chart-' + id), {
      type: 'line',
      data: {},
      options: options,
    });
    return r;
  }, {});

  // poll for chart data
  setInterval(poll, 5 * 60000); // 5m
  poll();
})();
