import ACTIONS from '../../src/constants/Constants';

import { getLanguage } from '../../src/actions/LanguageActionCreators.js';

describe('LanguageActionCreators', () => {

	it('should construct a fetch-action for the requested language', () => {

		const action = getLanguage('LANGUAGE');

		const expectedUrl = '/lang/LANGUAGE.json'; 

		const expectedAction = {
			[ACTIONS.CALL_API]: {
				types: [ACTIONS.GET_LANGUAGE, ACTIONS.GET_LANGUAGE_SUCCESS, ACTIONS.GET_LANGUAGE_ERROR],
				endpoint: expectedUrl,
				host: 'local'
			}
		};
		expect(action).toEqual(expectedAction);
	});
});