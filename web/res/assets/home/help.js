(() => {
  'use strict';

  // init popovers
  document.querySelectorAll('[data-bs-toggle="popover"]').forEach(
    (e) => new bootstrap.Popover(e)
  );
})();
