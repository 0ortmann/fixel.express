import React from 'react';
import path from 'path';
import Express from 'express';
import { renderToString } from 'react-dom/server';
import { Provider } from 'react-redux';
import { RouterContext, createMemoryHistory, match } from 'react-router';
import { syncHistoryWithStore } from 'react-router-redux';

import configureStore from '../src/store/configureStore.js';
import configureRoutes from '../src/routes.jsx';

const appConfig = require('../config/config.' + process.env.NODE_ENV + '.js');

const app = new Express();

const routes = configureRoutes(appConfig.routes);

if(process.env.NODE_ENV == 'development') {
	const webpack = require('webpack');

	const config = require(path.resolve(__dirname, '../webpack.config.' + process.env.NODE_ENV + '.js'));
	const compiler = webpack(config);
	app.use(require('webpack-dev-middleware')(compiler, {
		publicPath: config.output.publicPath,
		noInfo: true,
		serverSideRender: true
	}));

	app.use(require('webpack-hot-middleware')(compiler));
}

app.use(Express.static('public'));


function handleRender(req, res) {
	const memoryHistory = createMemoryHistory(req.url);
	let store = configureStore(memoryHistory, undefined, appConfig.api);
	const history = syncHistoryWithStore(memoryHistory, store);

	match({ history, routes, location: req.url }, (error, redirectLocation, renderProps) => {
		if (error) {
			res.status(500).send(error.message);
		}
		else if (redirectLocation) {
			res.redirect(302, redirectLocation.pathname + redirectLocation.search);
		}
		else if (renderProps) {
			fetchData(store, renderProps).then( () => {
				store = configureStore(memoryHistory, store.getState(), appConfig.api);
				const html = renderToString(
					<Provider store={store}>
						<RouterContext {...renderProps} />
					</Provider>
				);
				res.send(renderFullPage(html, store.getState()));
			});
		}
		else {
			res.status(404).send('Not found');
		}
	});
}

function fetchData(store, renderProps) {
	return new Promise((resolve) => {
		// delegate fetching the language inside app component (which is the root / wrapper comp)
		// it will always be accessable on index 0:
		let comp = renderProps.components[0].WrappedComponent;
		resolve(comp.fetchData(store));
	});
}

function renderFullPage(html, preloadedState) {
	return `
		<!DOCTYPE HTML>
		<html>

		<head>
			<meta charset='utf-8'>
			<meta name="viewport" content="width=device-width, initial-scale = 1.0, user-scalable=yes" />
			<meta name="description" content="Homepage of Felix Ortmann. I like to code, I love to think." />
			<title>Fixel Express</title>
			<link rel='stylesheet' type='text/css' href='/assets/bundle.css'>
			<link sizes="96x96" href="/favicon.ico" type="image/png" rel="icon">
		</head>

		<body>
			<div id='content'>${html}</div>
			<script>window.__PRELOADED_STATE__ = ${JSON.stringify(preloadedState)}</script>
			<script type='text/javascript' src='/assets/bundle.js' charset='utf-8'></script>
		</body>

		</html>
	`;
}

app.use(handleRender);

const server = app.listen(3300, function() {
	const host = server.address().address;
	const port = server.address().port;

	console.log('App listening at http://%s:%s', host, port);
});