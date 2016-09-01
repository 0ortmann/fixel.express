import keyMirror from 'key-mirror';

const ACTIONS = keyMirror({
	INSERT_TOKEN: null,
	INSERT_TOKEN_SUCCESS: null,
	INSERT_TOKEN_ERROR: null,
	NEW_GAME: null,
	NEW_GAME_SUCCESS: null,
	NEW_GAME_ERROR: null,
	CALL_API: null
});

export default ACTIONS;