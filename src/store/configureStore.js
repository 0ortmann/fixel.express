import configureApiMiddleware from '../middleware/Api.js';
import createLogger from 'redux-logger';
import lang from '../reducer/LanguageReducer.js';
import thunk from 'redux-thunk';
import { applyMiddleware, combineReducers, compose, createStore } from 'redux';
import { routerMiddleware, routerReducer } from 'react-router-redux';


const appConfig = require('../../config.' + process.env.NODE_ENV + '.js');

const reducer = combineReducers({
	routing: routerReducer,
	lang: lang
});

export default function configureStore(history, preloadedState, apiConfig) {
	const api = configureApiMiddleware(apiConfig.resourcesHost);
	let middleware = applyMiddleware(thunk, api, routerMiddleware(history));
	if (appConfig.logging) {
		middleware = applyMiddleware(thunk, api, createLogger(), routerMiddleware(history));
	}
	return createStore(
		reducer, 
		preloadedState, 
		compose(middleware)
	);
}