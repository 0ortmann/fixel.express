import 'isomorphic-fetch';
import ACTIONS from '../constants/Constants';


const API = 'http://localhost:5000';
const CALL_API = ACTIONS.CALL_API;

function getApi(endpoint) {
	const fullUrl = API + endpoint;

	return fetch( fullUrl )
		.then( res => {
			if (res.status >= 400) {
				return Promise.reject(res);
			}
			return res.json();
		});
}

function postApi(endpoint, body) {
	const fullUrl = API + endpoint;

	return fetch( fullUrl, {
		method: 'Post',
		body,
		})
		.then( res => {
			if (res.status >= 400) {
				return Promise.reject(res);
			}
			return res.json();
		});
}


function actionWith(origAction, newData) {
	const withoutCallApi = Object.assign({}, origAction, newData);
	delete withoutCallApi[CALL_API];
	return withoutCallApi;
}

export default store => next => action => {
	const apiCall = action[CALL_API];

	if (typeof apiCall === 'undefined') {
		return next(action);
	}
	const { endpoint, method, types } = apiCall;

	// mark as pending
	const [ requestType, successType, failureType ] = types;
	next( actionWith( action, { type: requestType }));
	if (method != 'Post') {
		
		return getApi(endpoint).then(
			response => next(actionWith( {
				response,
				type: successType
			})), 
			error => next(actionWith({
				type: failureType,
				error: error.message || 'Unable to get API'
			}))
		);
	}
	else if (method == 'Post') {
		return postApi(endpoint, body).then(
			response => next(actionWith( {
				response,
				type: successType
			})), 
			error => next(actionWith({
				type: failureType,
				error: error.message || 'Unable to post to API'
			}))
		);
	}
}
