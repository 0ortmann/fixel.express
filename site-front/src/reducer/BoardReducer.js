import ACTIONS from '../constants/Constants.js';

const initialState = {
	gameId: undefined,
	board: [[], [], [], [], [], [], []],
	isPlaying: false,
	playError: false
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
			board: apply(state.board, action.playerColumn, action.computerColumn),
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

function apply(board, pCol, cCol) {
	let muteBoard = {...board};
	muteBoard[pCol].push('player');
	muteBoard[cCol].push('computer');
	return muteBoard;
}