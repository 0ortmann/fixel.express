import 'babel-polyfill';

import React from 'react';
import { render } from 'react-dom';

import { Router, browserHistory } from 'react-router';
import { syncHistoryWithStore } from 'react-router-redux';
import { Provider } from 'react-redux';

import configureStore from './store/configureStore.js';
import configureRoutes from './routes.jsx';

const config = require('../config/config.' + process.env.NODE_ENV + '.js');

const store = configureStore(browserHistory, window.__PRELOADED_STATE__, config.api);
const history = syncHistoryWithStore(browserHistory, store);
const routes = configureRoutes();

render(
	<Provider store={store}>
		<Router routes={routes} history={history} />
	</Provider>,
	document.getElementById('content')
);
