import ACTIONS from '../constants/Constants.js';

const initialState =  {
	about: {},
	address: {},
	box: "",
	contact: {},
	fetchError: false,
	imprint: {},
	isFetching: false,
	privPolicy: {}
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
			address: action.response.address,
			box: action.response.box,
			contact: action.response.contact,
			fetchError: false,
			imprint: action.response.imprint,
			isFetching: false,
			privPolicy: action.response.privPolicy
		};
	default:
		return state;
	}
}