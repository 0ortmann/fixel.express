import ACTIONS from '../constants/Constants';

export function getLanguage(name) {
	return {
		[ACTIONS.CALL_API]: {
			method: 'Get',
			types: [ACTIONS.GET_LANGUAGE, ACTIONS.GET_LANGUAGE_SUCCESS, ACTIONS.GET_LANGUAGE_ERROR],
			endpoint: '/lang/' + name + '.json',
		}
	};
}