import ACTIONS from '../constants/Constants.js';

const initialState = {
	id: undefined,
	board: [[], [], [], [], [], [], []],
	isPlaying: false,
	playError: false,
	columns: 7,
	rows: 6,
	k: 4,
	level: 4,
	winner: '',
	isFetching: false
};

export default function game(state = initialState, action) {
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
			winner: action.response.winner,
			playError: false,
			isPlaying: false
		};
	case ACTIONS.NEW_GAME:
		return {
			...state,
			isFetching: true
		};
	case ACTIONS.NEW_GAME_ERROR:
		return {
			...state,
			id: undefined
		};
	case ACTIONS.NEW_GAME_SUCCESS: {
		let newBoard = [];
		action.response['board'].map(() => {
			newBoard.push([]);
		});
		return { 
			...state,
			id: action.response['id'],
			board: newBoard,
			isPlaying: initialState.isPlaying,
			playError: initialState.playError,
			columns: action.response['cols'],
			rows: action.response['rows'],
			k: action.response['win'],
			level: action.response['level'],

		};
	}
	default:
		return state;
	}
}

function apply(board, body, response) {
	const pCol = body.col; // player insert column
	const cCol = response.col; // computer insert column
	let muteBoard = [];
	board.map((row, i) => {
		muteBoard[i] = [...row];
	});
	muteBoard[pCol].push('player');
	if (cCol > -1) {
		muteBoard[cCol].push('computer');
	}
	return muteBoard;
}