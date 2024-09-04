(() => {
  'use strict';

  // cache edit dialog
  const EDIT_DIALOG = document.getElementById('edit-dialog');

  // populate edit dialog when shown
  EDIT_DIALOG.addEventListener('show.bs.modal', (ev) => {
    const data = ev.relatedTarget.dataset;

    document.getElementById('edit-id').value = data.id;
    document.getElementById('edit-name').value = data.name;
    document.getElementById('edit-color').value = data.color;
    document.getElementById('edit-sort').value = data.sort;
  });

  // bind to edit dialog save button click events
  document.getElementById('edit-save-btn').addEventListener('click', (ev) => {
    fetch('/api/home/current/edit', {
      method: 'POST',
      body: JSON.stringify({
        id: document.getElementById('edit-id').value,
        name: document.getElementById('edit-name').value,
        color: document.getElementById('edit-color').value,
        sort: +document.getElementById('edit-sort').value,
      }),
    }).then((r) => {
      if (!r.ok) {
        alert("Couldn't save changes");
        return;
      }

      // fire "saved" event (used by current and charts panels to
      // trigger refresh)
      EDIT_DIALOG.dispatchEvent(new CustomEvent('saved'));

      // dismiss dialog
      const close_btn_css = '#edit-dialog .modal-footer button[data-bs-dismiss]';
      document.querySelector(close_btn_css).click();
    });
  });
})();
