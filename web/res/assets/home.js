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
      <table>
        <thead>
          <tr>
            <th title='room' aria-label='room'>
              room
            </th>

            <th
              title='temperature (${unit.abbr})'
              aria-label='temperature (${unit.abbr})'
            >
              temp (&deg;${unit.abbr})
            </th>
          </tr>
        </thead>

        <tbody>${rows}</tbody>
      </table>
    `,

    row: (unit, {id, name, temp}) => `
      <tr id='${id}'>
        <td title='room' aria-label='room'>
          ${name}
        </td>

        <td
          style='text-align: right'
          title='temperature (${unit.abbr})'
          aria-label='temperature (${unit.abbr})'
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
