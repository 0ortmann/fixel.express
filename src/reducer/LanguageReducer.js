import ACTIONS from '../constants/Constants.js';

const initialState = {
	box: "",
	about: {},
	imprint: {},
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
			about: action.response.about,
			imprint: action.response.imprint,
			box: action.response.box,
			fetchError: false,
			isFetching: false
		};
	default:
		return state;
	}
}