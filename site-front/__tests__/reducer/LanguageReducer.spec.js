import ACTIONS from '../../src/constants/Constants';

import lang from '../../src/reducer/LanguageReducer.js';

describe('LanguageReducer', () => {

	const initialState =  {
		properties: undefined,
		isFetching: false,
		fetchError: false
	};

	it('should return the initial state', () => {
		expect(lang(undefined, {})).toEqual(initialState);
	});

	it('should handle "is-fetching" state', () => {
		const fetchAction = {
			type: ACTIONS.GET_LANGUAGE
		};
		expect(lang(undefined, fetchAction)).toEqual({
			properties: undefined,
			isFetching: true,
			fetchError: false
		});
	});

	it('should handle "fetch success"', () => {
		const response = {
			some: 'props',
			foo: 'bar'
		};
		const fetchAction = {
			type: ACTIONS.GET_LANGUAGE_SUCCESS,
			response
		};
		const previousState = {
			properties: undefined,
			isFetching: true,
			fetchError: false
		};
		const expectedNextState = {
			properties: response,
			isFetching: false,
			fetchError: false
		};

		expect(lang(previousState, fetchAction)).toEqual(expectedNextState);
	});

	it('should handle "fetch error"', () => {
		const fetchAction = {
			type: ACTIONS.GET_LANGUAGE_ERROR,
		};
		const previousState = {
			properties: { some: 'prop'},
			isFetching: true,
			fetchError: false
		};
		const expectedNextState = {
			properties: { some: 'prop'},
			isFetching: false,
			fetchError: true
		};

		expect(lang(previousState, fetchAction)).toEqual(expectedNextState);
	});
});