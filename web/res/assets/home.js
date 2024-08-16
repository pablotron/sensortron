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
    table: (unit, rows) => `
      <thead>
        <tr>
          <th
            title='Sensor location.'
            aria-label='Sensor location.'
          >
            Location
          </th>

          <th
            title='Temperature, measured in degrees ${unit.name}.'
            aria-label='Temperature, measured in degrees ${unit.name}.'
          >
            Temperature (&deg;${unit.abbr})
          </th>
        </tr>
      </thead>

      <tbody>${rows}</tbody>
    `,

    row: (unit, {id, name, temp}) => `
      <tr id='${id}'>
        <td
          title='Sensor location.'
          aria-label='Sensor location.'
        >
          ${h(name)}
        </td>

        <td
          style='text-align: right'
          title='Temperature of ${h(name)} (${unit.abbr}).'
          aria-label='Temperature of ${h(name)} (${unit.abbr}).'
        >
          ${(temp * unit.m + unit.b).toFixed(2)}
        </td>
      </tr>
    `,
  };

  // cache table wrapper
  const div = document.getElementById('current-temps');

  const poll = () => fetch('/api/poll', { method: 'POST' }).then(
    (r) => r.json()
  ).then((r) => {
    // get unit
    const unit = UNITS[document.querySelector('input.unit[type="radio"]:checked').value];
    div.innerHTML = T.table(unit, r.map((row) => T.row(unit, row)).join(''));
  });

  // bind click events
  document.querySelectorAll('input.unit, label.unit').forEach((e) => {
    e.addEventListener('click', () => { setTimeout(poll, 10) })
  });

  setInterval(poll, 10000);
  poll();
})();
