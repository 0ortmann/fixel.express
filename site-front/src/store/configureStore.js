import { createStore, applyMiddleware } from 'redux';
import thunk from 'redux-thunk';
import createLogger from 'redux-logger';
import api from '../middleware/Api.js';
import board from '../reducer/BoardReducer.js';


export default function configureStore(preloadedState) {
	return createStore(
		board, 
		preloadedState, 
		applyMiddleware(thunk, api, createLogger())
	);
}