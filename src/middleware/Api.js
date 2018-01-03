import 'isomorphic-fetch';
import ACTIONS from '../constants/Constants';

export function getApi(url) {
	return fetch( url )
		.then( res => {
			if (res.status >= 400) {
				return Promise.reject(res);
			}
			return res.json();
		});
}

export function actionWith(origAction, newData) {
	let withoutCallApi = Object.assign({}, origAction, newData);
	delete withoutCallApi[ACTIONS.CALL_API];
	return withoutCallApi;
}

export default function configureApiMiddleware(resourcesHost) {
	return () => next => action => {
		const apiCall = action[ACTIONS.CALL_API];

		if (typeof apiCall === 'undefined') {
			return next(action);
		}
		const { endpoint, method, types } = apiCall;

		// mark as pending
		const [ requestType, successType, failureType ] = types;
		next( actionWith( action, { type: requestType }));
		if (method == 'Get') {
			const url = resourcesHost + endpoint;
			return getApi(url)
				.then(response => next(actionWith({
					response,
					type: successType
				})))
				.catch(error => next(actionWith({
					type: failureType,
					error: error.message || 'Unable to get API'
				})));
		}
	};
}