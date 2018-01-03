import { createStore, applyMiddleware, compose, combineReducers } from 'redux';
import thunk from 'redux-thunk';
import createLogger from 'redux-logger';
import { routerReducer, routerMiddleware } from 'react-router-redux';

import configureApiMiddleware from '../middleware/Api.js';
import lang from '../reducer/LanguageReducer.js';

const appConfig = require('../../config.' + process.env.NODE_ENV + '.js');

const reducer = combineReducers({
	routing: routerReducer,
	lang: lang
});

export default function configureStore(history, preloadedState, apiConfig) {
	const api = configureApiMiddleware(apiConfig.resourcesHost, apiConfig.connectKHost);
	let middleware = applyMiddleware(thunk, api, routerMiddleware(history))
	if (!!appConfig.logging) {
		middleware = applyMiddleware(thunk, api, createLogger(), routerMiddleware(history))
	}
	return createStore(
		reducer, 
		preloadedState, 
		compose(middleware),
	);
}