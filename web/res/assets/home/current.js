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
    current_table: (unit, rows) => `
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

    current_row: (unit, {sensor, data}) => `
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
          title='Humidity in ${h(sensor.name)}.'
          aria-label='Humidity in ${h(sensor.name)}.'
        >
          ${(data.h * 100).toFixed(1)}%
        </td>
      </tr>
    `,
  };

  // cache current wrapper
  const current_el = document.getElementById('current');

  // poll for current sensor measurements
  const poll_current = () => fetch('/api/home/current/poll', { method: 'POST' }).then(
    (r) => r.json()
  ).then((r) => {
    const unit = UNITS[document.querySelector('input.unit[type="radio"]:checked').value],
          rows = r.map((row) => T.current_row(unit, row)).join('');
    current_el.innerHTML = T.current_table(unit, rows);
  });

  // bind to unit button click events
  document.querySelectorAll('input.unit, label.unit').forEach((e) => {
    e.addEventListener('click', () => { setTimeout(poll_current, 10) })
  });

  // bind to saved event
  document.getElementById('edit-dialog').addEventListener('saved', poll_current);

  // poll for current sensor measurements
  setInterval(poll_current, 10000); // 10s
  poll_current();
})();
