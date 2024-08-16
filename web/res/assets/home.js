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

  const poll = () => fetch('/api/home/poll', { method: 'POST' }).then(
    (r) => r.json()
  ).then((r) => {
    const unit = UNITS[document.querySelector('input.unit[type="radio"]:checked').value],
          rows = r.map((row) => T.current_row(unit, row)).join('');
    current_el.innerHTML = T.current_table(unit, rows);
  });

  const forecast = () => fetch('/api/home/forecast', { method: 'POST' }).then(
    (r) => r.json()
  ).then((r) => {
    const rows = r.properties.periods.slice(0, 6)
    forecast_el.innerHTML = rows.map((row) => T.forecast_row(row)).join('');
  });

  // bind click events
  document.querySelectorAll('input.unit, label.unit').forEach((e) => {
    e.addEventListener('click', () => { setTimeout(poll, 10) })
  });

  // bind to click events on names
  current_el.addEventListener('click', (ev) => {
    const data = ev.target.dataset;

    console.log(ev);
    ev.preventDefault();
    const name = prompt('Enter new name for sensor "' + data.name + '":', data.name);
    if (name !== null) {
      fetch('/api/home/edit', {
        method: 'POST',
        body: JSON.stringify({
          id: data.id,
          name: name,
        }),
      }).then((r) => poll());
    }
  }, true);

  setInterval(poll, 10000); // 10s
  poll();

  setInterval(forecast, 30 * 60000); // 30m
  forecast();
})();
