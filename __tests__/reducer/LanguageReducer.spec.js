import ACTIONS from '../../src/constants/Constants';

import lang from '../../src/reducer/LanguageReducer.js';

describe('LanguageReducer', () => {

	const initialState =  {
		about: {},
		address: {},
		contact: {},
		fetchError: false,
		footer: {},
		imprint: {},
		isFetching: false,
		privPolicy: {}
	};

	it('should return the initial state', () => {
		expect(lang(undefined, {})).toEqual(initialState);
	});

	it('should handle "is-fetching" state', () => {
		const fetchAction = {
			type: ACTIONS.GET_LANGUAGE
		};
		expect(lang(undefined, fetchAction)).toEqual({
			...initialState,
			isFetching: true,
		});
	});

	it('should handle "fetch success"', () => {
		const response = {
			about: {
				title: 'ABOUT'
			},
			imprint: {
				title: 'IMPRINT'
			},
		};
		const fetchAction = {
			type: ACTIONS.GET_LANGUAGE_SUCCESS,
			response
		};
		const previousState = {
			about: {},
			imprint: {},
			isFetching: true,
			fetchError: false
		};
		const expectedNextState = {
			about: {
				title: 'ABOUT'
			},
			imprint: {
				title: 'IMPRINT'
			},
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
			about: {},
			imprint: {},
			footer: {},
			isFetching: true,
			fetchError: false
		};
		const expectedNextState = {
			about: {},
			imprint: {},
			footer: {},
			isFetching: false,
			fetchError: true
		};

		expect(lang(previousState, fetchAction)).toEqual(expectedNextState);
	});
});