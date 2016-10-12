import ACTIONS from '../../src/constants/Constants';

import { newGame, insertToken } from '../../src/actions/ConnectKActionCreators.js';

describe('ConnectKActionCreators', () => {

	it('should construct a newGame fetch-action using all passed properties', () => {

		const properties = {
			foo: 'foo',
			bar: 'bar',
			batz: 1234
		};

		const action = newGame(properties);

		const expectedAction = {
			[ACTIONS.CALL_API]: {
				types: [ACTIONS.NEW_GAME, ACTIONS.NEW_GAME_SUCCESS, ACTIONS.NEW_GAME_ERROR],
				endpoint: '/new',
				properties
			}
		};
		expect(action).toEqual(expectedAction);
	});

	it('should construct a insertToken post-action for the passed "game" and "column"', () => {

		const game = 'SOME-UID-HERE';
		const column = 5;

		const action = insertToken(game, column);

		const expectedAction = {
			[ACTIONS.CALL_API]: {
				types: [ACTIONS.INSERT_TOKEN, ACTIONS.INSERT_TOKEN_SUCCESS, ACTIONS.INSERT_TOKEN_ERROR],
				endpoint: '/play',
				method: 'Post',
				body: { gameId: game, col: column }
			}
		};
		expect(action).toEqual(expectedAction);
	});
});