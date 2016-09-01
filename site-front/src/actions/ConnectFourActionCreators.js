import ACTIONS from '../constants/Constants';

export function insertToken(column) {
	return {
		[ACTIONS.CALL_API]: {
			types: [ACTIONS.INSERT_TOKEN, ACTIONS.INSERT_TOKEN_SUCCESS, ACTIONS.INSERT_TOKEN_ERROR],
			method: 'Post',
			endpoint: '/play',
			column
		}
	};
}

export function newGame() {
	return {
		[ACTIONS.CALL_API]: {
			types: [ACTIONS.NEW_GAME, ACTIONS.NEW_GAME_SUCCESS, ACTIONS.NEW_GAME_ERROR],
			endpoint: '/new'
		}
	};
}