import React, { Component } from 'react';
import { connect } from 'react-redux';

import ConnectKBoard from './ConnectKBoard.jsx';

import { newGame, insertToken } from '../../actions/ConnectFourActionCreators.js';


export class Game extends Component {
	
	constructor(props) {
		super(props);
		
		this.play = this.play.bind(this);
	}

	componentWillMount() {
		this.props.newGame();
	}

	play(column) {
		const { gameId, insertToken } = this.props;
		insertToken(gameId, column);
	}

	render() {
		const { board, cols, rows, newGame } = this.props;
		return (
			<main className='game'>
				<button onClick={newGame}>Ding Dong</button>
				<ConnectKBoard board={board} cols={cols} rows={rows} play={this.play} />
			</main>
		);
	}
}

function mapStateToProps(state, ownProps) {
	return {
		gameId: state.board.gameId,
		board: state.board.board,
		cols: state.board.columns,
		rows: state.board.rows
	}
}


export default connect(mapStateToProps, {
	newGame,
	insertToken
})(Game)