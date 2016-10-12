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
});