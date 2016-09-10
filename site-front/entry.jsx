import 'babel-polyfill';

import React from 'react';
import { render } from 'react-dom';

import { Router, browserHistory } from 'react-router';
import { syncHistoryWithStore } from 'react-router-redux';
import { Provider } from 'react-redux';

import configureStore from './src/store/configureStore.js';
import routes from './routes.js';

const store = configureStore(browserHistory, window.__PRELOADED_STATE__);
const history = syncHistoryWithStore(browserHistory, store);

render(
	<Provider store={store}>
		<Router routes={routes} history={history} />
	</Provider>,
	document.getElementById('content')
)
