import 'isomorphic-fetch';
import ACTIONS from '../constants/Constants';

const API = 'http://localhost:5000';
const CALL_API = ACTIONS.CALL_API;

function getApi(endpoint) {
	const fullUrl = endpoint.startsWith('http://')? endpoint : API + endpoint;

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
		method: 'POST',
		body: JSON.stringify(body),
		headers: {
			'Content-type': 'application/json'
		}
	})
	.then( res => {
		if (res.status >= 400) {
			return Promise.reject(res);
		}
		return res.json();
	});
}


function actionWith(origAction, newData) {
	let withoutCallApi = Object.assign({}, origAction, newData);
	delete withoutCallApi[CALL_API];
	return withoutCallApi;
}

function buildQueryString(obj) {
	let query = '?';
	for (const prop in obj) {
		query += prop + '=' + obj[prop] + '&';
	}
	return query;
}

export default () => next => action => {
	const apiCall = action[CALL_API];

	if (typeof apiCall === 'undefined') {
		return next(action);
	}
	const { endpoint, method, types, body, properties } = apiCall;

	// mark as pending
	const [ requestType, successType, failureType ] = types;
	next( actionWith( action, { type: requestType }));
	if (method != 'Post') {
		
		return getApi(endpoint + buildQueryString(properties)).then(
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
				body,
				response,
				type: successType
			})), 
			error => next(actionWith({
				type: failureType,
				error: error.message || 'Unable to post to API'
			}))
		);
	}
};