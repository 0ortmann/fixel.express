import { createStore, applyMiddleware, compose, combineReducers } from 'redux';
import thunk from 'redux-thunk';
import createLogger from 'redux-logger';
import { routerReducer, routerMiddleware } from 'react-router-redux';

import api from '../middleware/Api.js';
import game from '../reducer/GameReducer.js';

const reducer = combineReducers({
	routing: routerReducer,
	game: game
});


export default function configureStore(history, preloadedState) {
	return createStore(
		reducer, 
		preloadedState, 
		compose( // need me?
			applyMiddleware(thunk, api, createLogger(), routerMiddleware(history))
		)
	);
}