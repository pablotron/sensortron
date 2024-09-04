(() => {
  'use strict';

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
    table: (unit, rows) => `
      <thead>
        <tr>
          <th
            title='Sensor name.'
            aria-label='Sensor name.'
          >
            Name
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

    row: (unit, {sensor, data}) => `
      <tr id='${h(sensor.id)}'>
        <td
          title='Sensor name.'
          aria-label='Sensor name.'
        >
          <a
            href='#'

            title='Sensor name.'
            aria-label='Sensor name.'

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
          title='${h(sensor.name)} humidity (%).'
          aria-label='${h(sensor.name)} (%).'
        >
          ${(data.h * 100).toFixed(1)}%
        </td>
      </tr>
    `,
  };

  // poll for current sensor measurements
  const poll = () => fetch('/api/home/current/poll', { method: 'POST' }).then(
    (r) => r.json()
  ).then((r) => {
    const unit = UNITS[document.querySelector('input.unit[type="radio"]:checked').value],
          rows = r.map((row) => T.row(unit, row)).join('');
    document.getElementById('current').innerHTML = T.table(unit, rows);
  });

  // bind to unit button click events
  document.querySelectorAll('input.unit, label.unit').forEach((e) => {
    e.addEventListener('click', () => { setTimeout(poll, 10) })
  });

  // bind to saved event
  document.getElementById('edit-dialog').addEventListener('saved', poll);

  // poll for current sensor measurements
  setInterval(poll, 10000); // 10s
  poll();
})();
