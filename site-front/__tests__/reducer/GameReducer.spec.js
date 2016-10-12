import ACTIONS from '../../src/constants/Constants';

import game from '../../src/reducer/GameReducer.js';

describe('GameReducer new game', () => {

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

	it('should return the initial state', () => {
		expect(game(undefined, {})).toEqual(initialState);
	});

	it('should handle "is-fetching" state for new-game action', () => {
		const fetchNewGameAction = {
			type: ACTIONS.NEW_GAME
		};
		expect(game(undefined, fetchNewGameAction)).toEqual({
			...initialState,
			isFetching: true
		});
	});

	it('should handle "new game success"', () => {
		const response = {
			// connect-6 on 10x10 board, level 8
			cols: 10,
			rows: 10,
			level: 8,
			win: 6,
			board: [null,null,null,null,null,null,null,null,null,null],
			id: 'SOME-UID-HERE'
		};
		const fetchAction = {
			type: ACTIONS.NEW_GAME_SUCCESS,
			response
		};
		const previousState = {
			id: 'SOME-OTHER-UID-HERE',
			board: [[], [], [], [], [], [], [], []],
			isPlaying: false,
			playError: false,
			columns: 8,
			rows: 8,
			k: 4,
			level: 4,
			winner: '',
			isFetching: false
		};
		const expectedNextState = {
			id: 'SOME-UID-HERE',
			board: [[], [], [], [], [], [], [], [], [], []],
			isPlaying: false,
			playError: false,
			columns: 10,
			rows: 10,
			k: 6,
			level: 8,
			winner: '',
			isFetching: false
		};

		expect(game(previousState, fetchAction)).toEqual(expectedNextState);
	});

	it('should handle "new game error"', () => {
		const fetchAction = {
			type: ACTIONS.NEW_GAME_ERROR,
		};
		const previousState = {
			id: 'SOME-OTHER-UID-HERE',
			board: [[], [], [], [], [], [], [], []],
			isPlaying: false,
			playError: false,
			columns: 8,
			rows: 8,
			k: 4,
			level: 4,
			winner: '',
			isFetching: false
		};
		const expectedNextState = {
			...previousState,
			id: undefined
		};

		expect(game(previousState, fetchAction)).toEqual(expectedNextState);
	});

	it('should handle "is-playing" state for insert token action', () => {
		const insertTokenAction = {
			type: ACTIONS.INSERT_TOKEN
		};
		expect(game(undefined, insertTokenAction)).toEqual({
			...initialState,
			isPlaying: true
		});
	});

	it('should handle apply tokens accordingly on "insert token success"', () => {
		const response = {
			col: 5,
			winner: ''
		};
		const insertAction = {
			type: ACTIONS.INSERT_TOKEN_SUCCESS,
			response, // compuer did this
			body: { col: 1 } // player did this
		};
		const previousState = {
			id: 'SOME-UID-HERE',
			board: [[], [], [], [], [], [], [], []],
			isPlaying: true,
			playError: false,
			columns: 8,
			rows: 8,
			k: 4,
			level: 4,
			winner: '',
			isFetching: false
		};
		const expectedNextState = {
			...previousState,
			board: [[], ['player'], [], [], [], ['computer'], [], []],
			isPlaying: false,
		};

		expect(game(previousState, insertAction)).toEqual(expectedNextState);
	});

	it('should set the winner correctly on "insert token success"', () => {
		const response = {
			col: -1,
			winner: 'player'
		};
		const insertAction = {
			type: ACTIONS.INSERT_TOKEN_SUCCESS,
			response,
			body: { col: 1 }
		};
		const previousState = {
			id: 'SOME-UID-HERE',
			board: [[], [], [], [], [], [], [], []],
			isPlaying: true,
			playError: false,
			columns: 8,
			rows: 8,
			k: 4,
			level: 4,
			winner: '',
			isFetching: false
		};
		const expectedNextState = {
			...previousState,
			id: 'SOME-UID-HERE',
			board: [[], ['player'], [], [], [], [], [], []],
			isPlaying: false,
			winner: 'player',
		};

		expect(game(previousState, insertAction)).toEqual(expectedNextState);
	});

	it('should handle "new game error"', () => {
		const fetchAction = {
			type: ACTIONS.INSERT_TOKEN_ERROR,
		};
		const previousState = {
			id: 'SOME-OTHER-UID-HERE',
			board: [[], [], [], [], [], [], [], []],
			isPlaying: false,
			playError: false,
			columns: 8,
			rows: 8,
			k: 4,
			level: 4,
			winner: '',
			isFetching: false
		};
		const expectedNextState = {
			...previousState,
			isPlaying: false,
			playError: true
		};

		expect(game(previousState, fetchAction)).toEqual(expectedNextState);
	});
});
