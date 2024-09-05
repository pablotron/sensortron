#!/bin/bash

#
# install-assets.sh: Copy CSS and JS assets from the node_modules/
# directory to res/assets/ directory.
#

set -eu

# copy bootstrap assets
mkdir -p res/assets/bootstrap-5.3.3/{css,js}
cp node_modules/bootstrap/dist/js/bootstrap.bundle.min.js{,.map} res/assets/bootstrap-5.3.3/js/
cp node_modules/bootstrap/dist/css/bootstrap.min.css{,.map} res/assets/bootstrap-5.3.3/css/

# copy chartjs and chartjs-adapter-date-fns
cp node_modules/chart.js/dist/chart.umd.js res/assets/chart-4.4.3.min.js
cp node_modules/chartjs-adapter-date-fns/dist/chartjs-adapter-date-fns.bundle.min.js res/assets/chartjs-adapter-date-fns-3.0.0.bundle.min.js
