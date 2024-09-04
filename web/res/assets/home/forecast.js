(() => {
  'use strict';

  // templates
  const T = {
    row: (row) => `
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
  };

  // cache forecast wrapper
  const forecast_el = document.getElementById('forecast');

  // poll for current forecast
  const poll = () => fetch('/api/home/forecast/poll', { method: 'POST' }).then(
    (r) => r.json()
  ).then((r) => {
    forecast_el.dataset.forecast = JSON.stringify(r);
    const rows = r.properties.periods.slice(0, 8);
    forecast_el.innerHTML = rows.map((row) => T.row(row)).join('');
  });

  // start polling
  setInterval(poll, 30 * 60000); // 30m
  poll();
})();
