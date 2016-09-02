import ACTIONS from '../constants/Constants.js';

const initialState = {
	gameId: undefined,
	board: [[], [], [], [], [], [], []],
	isPlaying: false,
	playError: false,
	columns: 7,
	rows: 6
};

export default function board(state = initialState, action) {
	switch(action.type) {
	case ACTIONS.INSERT_TOKEN:
		return {
			...state,
			isPlaying: true
		};
	case ACTIONS.INSERT_TOKEN_ERROR:
		return {
			...state,
			playError: true,
			isPlaying: false
		};
	case ACTIONS.INSERT_TOKEN_SUCCESS:
		return { 
			...state, 
			board: apply(state.board, action.body, action.response),
			playError: false,
			isPlaying: false
		};
	case ACTIONS.NEW_GAME_ERROR:
		return {
			...state,
			gameId: undefined
		};
	case ACTIONS.NEW_GAME_SUCCESS:
		return { 
			...state, 
			gameId: action.response['Id']
		};
	default:
		return state;
	}
}

function apply(board, body, response) {
	const pCol = body.col; // player insert column
	const cCol = response.col; // computer insert column
	let muteBoard = [...board];
	muteBoard[pCol].push('player');
	muteBoard[cCol].push('computer');
	return muteBoard;
}