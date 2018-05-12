import ACTIONS from '../../src/constants/Constants';

import lang from '../../src/reducer/LanguageReducer.js';

describe('LanguageReducer', () => {

	const initialState =  {
		box: "",
		about: {},
		imprint: {},
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
			box: "",
			about: {},
			imprint: {},
			isFetching: true,
			fetchError: false
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
			box: 'BOX',
		};
		const fetchAction = {
			type: ACTIONS.GET_LANGUAGE_SUCCESS,
			response
		};
		const previousState = {
			about: {},
			imprint: {},
			box: '',
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
			box: 'BOX',
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
			box: '',
			isFetching: true,
			fetchError: false
		};
		const expectedNextState = {
			about: {},
			imprint: {},
			box: '',
			isFetching: false,
			fetchError: true
		};

		expect(lang(previousState, fetchAction)).toEqual(expectedNextState);
	});
});