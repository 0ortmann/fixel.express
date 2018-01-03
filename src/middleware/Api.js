import 'isomorphic-fetch';
import ACTIONS from '../constants/Constants';

const CALL_API = ACTIONS.CALL_API;

export function getApi(url) {
	return fetch( url )
		.then( res => {
			if (res.status >= 400) {
				return Promise.reject(res);
			}
			return res.json();
		});
}

export function postApi(url, body) {
	return fetch( url, {
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


export function actionWith(origAction, newData) {
	let withoutCallApi = Object.assign({}, origAction, newData);
	delete withoutCallApi[CALL_API];
	return withoutCallApi;
}

export function buildQueryString(obj) {
	let query = '?';
	for (const prop in obj) {
		query += prop + '=' + obj[prop] + '&';
	}
	return query;
}

export function buildUrl(resourcesHost, endpoint, host) {
	resourcesHost = resourcesHost.endsWith('/')? resourcesHost.substring(0, resourcesHost.length -1) : resourcesHost;
	return resourcesHost + endpoint;
}

export default function configureApiMiddleware(resourcesHost) {
	return () => next => action => {
		const apiCall = action[CALL_API];

		if (typeof apiCall === 'undefined') {
			return next(action);
		}
		const { endpoint, method, types, body, properties, host } = apiCall;

		// mark as pending
		const [ requestType, successType, failureType ] = types;
		next( actionWith( action, { type: requestType }));
		if (method != 'Post') {
			const url = buildUrl(resourcesHost, endpoint, host);
			return getApi(url + buildQueryString(properties)).then(
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
			const url = buildUrl(resourcesHost, endpoint, host);
			return postApi(url, body).then(
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
}