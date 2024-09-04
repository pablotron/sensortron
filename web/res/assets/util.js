'use strict';

// html escape (replaceall explicit)
const h = (v) => {
  return v.toString().replaceAll('&', '&amp;')
    .replaceAll('<', '&lt;')
    .replaceAll('>', '&gt;')
    .replaceAll("'", '&apos;')
    .replaceAll('"', '&quot;');
};
