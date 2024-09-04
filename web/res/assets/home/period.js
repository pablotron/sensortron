(() => {
  'use strict';

  // templates
  const T = {
    body: (row) => `
      <img
        src='${h(row.icon)}'
        class='rounded float-start me-2'
        title='${h(row.name)}: ${h(row.detailedForecast)}'
        aria-label='${h(row.name)}: ${h(row.detailedForecast)}'
        alt='${h(row.name)}: ${h(row.detailedForecast)}'
      />

      <dl class='row'>
        <dt
          class='col-sm-3'
          title='Forecast time period'
          aria-label='Forecast time period'
        >
          Time
        </dt>
        <dd
          class='col-sm-9'
          title='Forecast time period'
          aria-label='Forecast time period'
        >
          ${h(row.start)} - ${h(row.end)}
        </dd>

        <dt class='col-sm-3' title='Temperature' aria-label='Temperature'>
          Temperature
        </dt>
        <dd class='col-sm-9' title='Temperature' aria-label='Temperature'>
          ${row.temperature}&deg;F
        </dd>

        <dt
          class='col-sm-3'
          title='Probability of precipitation'
          aria-label='Probability of precipitation'
        >
          Precipitation
        </dt>
        <dd
          class='col-sm-9'
          title='Probability of precipitation'
          aria-label='Probability of precipitation'
        >
          ${row.probabilityOfPrecipitation.value ?? '0'}%
        </dd>

        <dt
          class='col-sm-3'
          title='Wind speed and direction'
          aria-label='Wind speed and direction'
        >
          Wind
        </dt>
        <dd
          class='col-sm-9'
          title='Wind speed and direction'
          aria-label='Wind speed and direction'
        >
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

  // populate period dialog when shown
  document.getElementById('period-dialog').addEventListener('show.bs.modal', (ev) => {
    // parse row, format dates
    const row = JSON.parse(ev.relatedTarget.dataset.row);
    row.start = (new Date(row.startTime)).toLocaleString();
    row.end = (new Date(row.endTime)).toLocaleString();

    // set modal title and body
    document.getElementById('period-dialog-title').textContent = row.name;
    document.getElementById('period-dialog-body').innerHTML = T.body(row);
  });
})();
