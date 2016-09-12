import ACTIONS from '../constants/Constants.js';

const initialState = {
	properties: undefined,
	isFetching: false,
	fetchError: false
};

export default function lang(state = initialState, action) {
	switch(action.type) {
	case ACTIONS.GET_LANGUAGE:
		return {
			...state,
			isFetching: true
		};
	case ACTIONS.GET_LANGUAGE_ERROR:
		return {
			...state,
			fetchError: true,
			isFetching: false
		};
	case ACTIONS.GET_LANGUAGE_SUCCESS:
		return { 
			...state,
			properties: action.response,
			fetchError: false,
			isFetching: false
		};
	default:
		return state;
	}
}