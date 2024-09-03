(() => {
  'use strict';

  // html escape (replaceall explicit)
  const h = (v) => {
    return v.toString().replaceAll('&', '&amp;')
      .replaceAll('<', '&lt;')
      .replaceAll('>', '&gt;')
      .replaceAll("'", '&apos;')
      .replaceAll('"', '&quot;');
  };

  // temperature units
  const UNITS = {
    f: {
      abbr: 'F',
      name: 'Fahrenheit',
      m: 1.8,
      b: 32,
    },

    c: {
      abbr: 'C',
      name: 'Celcius',
      m: 1.0,
      b: 0,
    },
  };

  // chart options
  const CHART_OPTIONS = [{
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

  // templates
  const T = {
    current_table: (unit, rows) => `
      <thead>
        <tr>
          <th
            title='Sensor location.'
            aria-label='Sensor location.'
          >
            Location
          </th>

          <th
            style='text-align: right'
            title='Temperature (in ${unit.name}).'
            aria-label='Temperature (in ${unit.name}).'
          >
            Temperature (&deg;${unit.abbr})
          </th>

          <th
            style='text-align: right'
            title='Relative humidity (percent).'
            aria-label='Relative humidity (percent).'
          >
            Humidity (%)
          </th>
        </tr>
      </thead>

      <tbody>${rows}</tbody>
    `,

    current_row: (unit, {sensor, data}) => `
      <tr id='${h(sensor.id)}'>
        <td
          title='Sensor location.'
          aria-label='Sensor location.'
        >
          <a
            href='#'

            title='Sensor location.'
            aria-label='Sensor location.'

            data-bs-toggle='modal'
            data-bs-target='#edit-dialog'

            data-id='${h(sensor.id)}'
            data-name='${h(sensor.name)}'
            data-color='${h(sensor.color)}'
            data-sort='${h(sensor.sort)}'
          >
            ${h(sensor.name)}
          </a>
        </td>

        <td
          style='text-align: right'
          title='${h(sensor.name)} temperature (${unit.abbr}).'
          aria-label='${h(sensor.name)} temperature (${unit.abbr}).'
        >
          ${(data.t * unit.m + unit.b).toFixed(2)}
        </td>

        <td
          style='text-align: right'
          title='Humidity in ${h(sensor.name)}.'
          aria-label='Humidity in ${h(sensor.name)}.'
        >
          ${(data.h * 100).toFixed(1)}%
        </td>
      </tr>
    `,

    forecast_row: (row) => `
      <li
        class='list-group-item'
        title='${h(row.name)}: ${h(row.detailedForecast)}'
        aria-label='${h(row.name)}: ${h(row.detailedForecast)}'
        data-row='${h(JSON.stringify(row))}'
        data-bs-toggle='modal'
        data-bs-target='#period-dialog'
      >
        <img
          src='${h(row.icon)}'
          class='rounded float-start me-2'
          title='${h(row.name)}: ${h(row.detailedForecast)}'
          aria-label='${h(row.name)}: ${h(row.detailedForecast)}'
          alt='${h(row.name)}: ${h(row.detailedForecast)}'
        />

        <h5>${h(row.name)}</h5>
        Temperature: <b>${row.temperature}&deg;F</b>,
        Rain: <b>${row.probabilityOfPrecipitation.value ?? '0'}%</b><br/>
        ${h(row.shortForecast)}
      </li>
    `,

    period_dialog_body: (row) => `
      <img
        src='${h(row.icon)}'
        class='rounded float-start me-2'
        title='${h(row.name)}: ${h(row.detailedForecast)}'
        aria-label='${h(row.name)}: ${h(row.detailedForecast)}'
        alt='${h(row.name)}: ${h(row.detailedForecast)}'
      />

      <dl class='row'>
        <dt class='col-sm-3' title='Time' aria-label='Time'>
          Time
        </dt>
        <dd class='col-sm-9' title='Time' aria-label='Time'>
          ${h(row.start)} - ${h(row.end)}
        </dd>

        <dt class='col-sm-3' title='Temperature' aria-label='Temperature'>
          Temperature
        </dt>
        <dd class='col-sm-9' title='Temperature' aria-label='Temperature'>
          ${row.temperature}&deg;F
        </dd>

        <dt class='col-sm-3' title='Precipitation' aria-label='Precipitation'>
          Precipitation
        </dt>
        <dd class='col-sm-9' title='Precipitation' aria-label='Precipitation'>
          ${row.probabilityOfPrecipitation.value ?? '0'}%
        </dd>

        <dt class='col-sm-3' title='Wind' aria-label='Wind'>
          Wind
        </dt>
        <dd class='col-sm-9' title='Wind' aria-label='Wind'>
          ${row.windSpeed} ${row.windDirection}
        </dd>

        <dt class='col-sm-3' title='Forecast' aria-label='Forecast'>
          Forecast
        </dt>
        <dd class='col-sm-9' title='Forecast' aria-label='Forecast'>
          ${h(row.detailedForecast)}
        </dd>
      </dl>
    `,
  };

  // cache current data wrapper and forecast wrapper
  const current_el = document.getElementById('current'),
        forecast_el = document.getElementById('forecast');

  // poll for current sensor measurements
  const poll_current = () => fetch('/api/home/current/poll', { method: 'POST' }).then(
    (r) => r.json()
  ).then((r) => {
    const unit = UNITS[document.querySelector('input.unit[type="radio"]:checked').value],
          rows = r.map((row) => T.current_row(unit, row)).join('');
    current_el.innerHTML = T.current_table(unit, rows);
  });

  // poll for current forecast
  const poll_forecast = () => fetch('/api/home/forecast/poll', { method: 'POST' }).then(
    (r) => r.json()
  ).then((r) => {
    forecast_el.dataset.forecast = JSON.stringify(r);
    const rows = r.properties.periods.slice(0, 8);
    forecast_el.innerHTML = rows.map((row) => T.forecast_row(row)).join('');
  });

  // poll for current chart data
  const poll_charts = () => fetch('/api/home/charts/poll', { method: 'POST' }).then(
    (r) => r.json()
  ).then((r) => {
    // get hours, convert to slice start
    const start = 4 * -document.querySelector('input.time[type="radio"]:checked').value;
    for (let k of Object.keys(r)) {
      if (k in charts) {
        {
          // filter chart data based on time filter
          const labels = r[k].labels.slice(start),
                datasets = r[k].datasets.map((set) => {
                  set.data = set.data.slice(start);
                  return set;
                });
          r[k].labels = labels;
          r[k].datasets = datasets;
        }

        charts[k].data = r[k];
        charts[k].update();
      }
    }
  });

  // bind to unit button click events
  document.querySelectorAll('input.unit, label.unit').forEach((e) => {
    e.addEventListener('click', () => { setTimeout(poll_current, 10) })
  });

  // populate edit dialog when shown
  document.getElementById('edit-dialog').addEventListener('show.bs.modal', (ev) => {
    const data = ev.relatedTarget.dataset;

    document.getElementById('edit-id').value = data.id;
    document.getElementById('edit-name').value = data.name;
    document.getElementById('edit-color').value = data.color;
    document.getElementById('edit-sort').value = data.sort;
  });

  // bind to edit dialog save button click events
  document.getElementById('edit-save-btn').addEventListener('click', (ev) => {
    fetch('/api/home/current/edit', {
      method: 'POST',
      body: JSON.stringify({
        id: document.getElementById('edit-id').value,
        name: document.getElementById('edit-name').value,
        color: document.getElementById('edit-color').value,
        sort: +document.getElementById('edit-sort').value,
      }),
    }).then((r) => {
      if (!r.ok) {
        alert("Couldn't save changes");
        return;
      }

      // refresh current and charts
      poll_current();
      poll_charts();

      // dismiss dialog
      const close_btn_css = '#edit-dialog .modal-footer button[data-bs-dismiss]';
      document.querySelector(close_btn_css).click();
    });
  });

  // populate period dialog when shown
  document.getElementById('period-dialog').addEventListener('show.bs.modal', (ev) => {
    // parse row, format dates
    const row = JSON.parse(ev.relatedTarget.dataset.row);
    row.start = (new Date(row.startTime)).toLocaleString();
    row.end = (new Date(row.endTime)).toLocaleString();

    // render modal title and body
    document.getElementById('period-dialog-title').textContent = row.name;
    document.getElementById('period-dialog-body').innerHTML = T.period_dialog_body(row);
  });

  // bind click events on time filter
  document.querySelectorAll('input.time, label.time').forEach((e) => {
    e.addEventListener('click', () => { setTimeout(poll_charts, 10) })
  });

  // poll for current sensor measurements
  setInterval(poll_current, 10000); // 10s
  poll_current();

  // poll for current forecast
  setInterval(poll_forecast, 30 * 60000); // 30m
  poll_forecast();

  // init charts
  Chart.defaults.color = '#eee';
  const charts = CHART_OPTIONS.reduce((r, {id, options}) => {
    r[id] = new Chart(document.getElementById('chart-' + id), {
      type: 'line',
      data: {},
      options: options,
    });
    return r;
  }, {});

  // poll for chart data
  setInterval(poll_charts, 5 * 60000); // 5m
  poll_charts();
})();
