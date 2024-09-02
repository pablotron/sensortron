(() => {
  'use strict';

  // html escape (replaceall explicit)
  const h = (v) => {
    return v.replaceAll('&', '&amp;')
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

    current_row: (unit, {id, name, data}) => `
      <tr id='${id}'>
        <td
          title='Sensor location.'
          aria-label='Sensor location.'
        >
          <a
            href='#'
            data-id='${id}'
            data-name='${h(name)}'
            title='Sensor location.'
            aria-label='Sensor location.'
          >
            ${h(name)}
          </a>
        </td>

        <td
          style='text-align: right'
          title='${h(name)} temperature (${unit.abbr}).'
          aria-label='${h(name)} temperature (${unit.abbr}).'
        >
          ${(data.t * unit.m + unit.b).toFixed(2)}
        </td>

        <td
          style='text-align: right'
          title='Humidity in ${h(name)}.'
          aria-label='Humidity in ${h(name)}.'
        >
          ${(data.h * 100).toFixed(1)}%
        </td>
      </tr>
    `,

    forecast_row: ({name, icon, shortForecast, detailedForecast, temperature, windSpeed, windDirection}) => `
      <li class='list-group-item'>
        <img
          src='${icon}'
          class='rounded float-start me-2'
          title='$(h(name)): ${h(shortForecast)}'
          aria-label='$(h(name)): ${h(shortForecast)}'
          alt='$(h(name)): ${h(shortForecast)}'
        />

        <h5>${h(name)}</h5>
        Temperature: <b>${temperature} F</b>, Wind: ${windSpeed} ${windDirection}<br/>
        ${h(shortForecast)}
      </li>
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
    const rows = r.properties.periods.slice(0, 4);
    forecast_el.innerHTML = rows.map((row) => T.forecast_row(row)).join('');
  });

  // poll for current chart data
  const poll_charts = () => fetch('/api/home/charts/poll', { method: 'POST' }).then(
    (r) => r.json()
  ).then((r) => {
    for (let k of Object.keys(r)) {
      if (k in charts) {
        charts[k].data = r[k];
        charts[k].update();
      }
    }
  });

  // bind click events
  document.querySelectorAll('input.unit, label.unit').forEach((e) => {
    e.addEventListener('click', () => { setTimeout(poll, 10) })
  });

  // bind to click events on names
  current_el.addEventListener('click', (ev) => {
    if (ev.target.tagName === 'A') {
      const data = ev.target.dataset;

      console.log(ev);
      ev.preventDefault();
      const name = prompt('Enter new name for sensor "' + data.name + '":', data.name);
      if (name !== null) {
        fetch('/api/home/current/edit', {
          method: 'POST',
          body: JSON.stringify({
            id: data.id,
            name: name,
          }),
        }).then((r) => poll());
      }
    }
  }, true);

  // poll for current sensor measurements
  setInterval(poll_current, 10000); // 10s
  poll_current();

  // poll for current forecast
  setInterval(poll_forecast, 30 * 60000); // 30m
  poll_forecast();

  // init chart
  const charts = {
    t: new Chart(document.getElementById('chart-t'), {
      type: 'line',
      data: {},
      options: {
        maintainAspectRatio: false,
      }
    }),
  };

  // poll for chart data
  setInterval(poll_charts, 5 * 60000); // 5m
  poll_charts();
})();
