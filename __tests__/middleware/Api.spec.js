import fetchMock from 'fetch-mock';

import ACTIONS from '../../src/constants/Constants';

import configureApiMiddleware, {
	getApi, postApi, actionWith, buildQueryString, buildUrl
} from '../../src/middleware/Api.js';

const resourcesEndpoint = 'http://RESOURCES';

const dispatchWith = (api, action) => {
	let dispatched = null;
	// "mock" of the dispatch function
	const dispatch = api()(actionAttempt => dispatched = actionAttempt);
	dispatch(action);
	return dispatched;
};

const awaitDispatch = (api, action) => {
	// "mock" of the dispatch function
	const dispatch = api()(() => {});
	return dispatch(action);
};

describe('API middleware', () => {
	
	const api = configureApiMiddleware(resourcesEndpoint);

	it('should dispatch any kind of action', () => {
		const action = {
			type: 'SOME_UNKNOWN_TYPE',
			property: 'foobar'
		};

		expect(dispatchWith(api, action)).toEqual(action);
	});

	it('should dispatch fetching state for GET "CALL_API" action that is not resolved so far', () => {
		const action = {
			[ACTIONS.CALL_API]: {
				types: ['NEW', 'SUCCESS', 'ERROR'],
				endpoint: '/succ',
				other: 'some property'
			}
		};
		const expectedDispatch = {
			type: 'NEW'
		};
		expect(dispatchWith(api, action)).toEqual(expectedDispatch);
	});

	it('should dispatch fetching state for POST "CALL_API" action that is not resolved so far', () => {
		const action = {
			[ACTIONS.CALL_API]: {
				types: ['NEW', 'SUCCESS', 'ERROR'],
				endpoint: '/succ',
				method: 'post',
				other: 'some property'
			}
		};
		const expectedDispatch = {
			type: 'NEW'
		};
		expect(dispatchWith(api, action)).toEqual(expectedDispatch);
	});

	it('should dispatch success for GET "CALL_API" action that resolved', () => {
		const action = {
			[ACTIONS.CALL_API]: {
				types: ['NEW', 'SUCCESS', 'ERROR'],
				endpoint: '/succ',
				other: 'some property'
			}
		};
		const expectedSuccessDispatch = {
			type: 'SUCCESS',
			other: 'some property'
		};
		
		awaitDispatch(api, action).then(result => expect(result).toEqual(expectedSuccessDispatch));
	});

	it('should dispatch success for POST "CALL_API" action that resolved', () => {
		const action = {
			[ACTIONS.CALL_API]: {
				types: ['NEW', 'SUCCESS', 'ERROR'],
				endpoint: '/succ',
				method: 'post',
				other: 'some property'
			}
		};
		const expectedSuccessDispatch = {
			type: 'SUCCESS',
			other: 'some property'
		};
		
		awaitDispatch(api, action).then(result => expect(result).toEqual(expectedSuccessDispatch));
	});

	it('should dispatch error for GET "CALL_API" action that rejects', () => {
		const action = {
			[ACTIONS.CALL_API]: {
				types: ['NEW', 'SUCCESS', 'ERROR'],
				endpoint: '/notfound',
				other: 'some property'
			}
		};
		const expectedErrorDispatch = {
			type: 'ERROR',
			other: 'some property'
		};
		
		awaitDispatch(api, action).then(result => expect(result).toEqual(expectedErrorDispatch));
	});

	it('should dispatch error for POST "CALL_API" action that rejects', () => {
		const action = {
			[ACTIONS.CALL_API]: {
				types: ['NEW', 'SUCCESS', 'ERROR'],
				endpoint: '/notfound',
				method: 'post',
				other: 'some property'
			}
		};
		const expectedErrorDispatch = {
			type: 'ERROR',
			other: 'some property'
		};
		
		awaitDispatch(api, action).then(result => expect(result).toEqual(expectedErrorDispatch));
	});
	
});

describe('API query / url builder', () => {

	it('should construct a valid Query string using all provided properties', () => {
		const props = { c: 1, r: 5, l: 'super-duper', k: 'foobar' };
		const qs = buildQueryString(props);
		expect(qs).toEqual('?c=1&r=5&l=super-duper&k=foobar&'); 
	});

	it('should be taken care of double slashes during URL construction', () => {
		const endpoint = '/endpoint?foo=bar';
		const url = buildUrl(resourcesEndpoint + '/', endpoint);
		expect(url).toEqual(resourcesEndpoint + endpoint);
	});
});

describe('API calls', () => {

	it('should reject on GET requests with erroneous stati', () => {
		return getApi(resourcesEndpoint + '/notfound').then(
			() => { throw 'Reached then branch, but should have failed'; }, 
			() => {}
		);
	});

	it('should reject on POST requests with erroneous stati', () => {
		return postApi(resourcesEndpoint + '/notfound', { post: 'body' }).then(
			() => { throw 'Reached then branch, but should have failed'; }, 
			() => {}
		);
	});

	it('should return server response on GET requests', () => {
		return getApi(resourcesEndpoint + '/succ').then(
			res => { expect(res).toEqual({ get: 'success' }); }, 
			() => { throw 'Request should not have failed'; }
		);
	});

	it('should return server response on GET requests', () => {
		return postApi(resourcesEndpoint + '/succ', { post: 'body' }).then(
			res => { expect(res).toEqual({ post: 'success' }); }, 
			() => { throw 'Request should not have failed'; }
		);
	});

});