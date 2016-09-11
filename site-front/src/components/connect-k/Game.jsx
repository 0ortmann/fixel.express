import React, { Component } from 'react';
import { connect } from 'react-redux';

import ConnectKBoard from './ConnectKBoard.jsx';
import PropertyPanel from './PropertyPanel.jsx';

import { newGame, insertToken } from '../../actions/ConnectFourActionCreators.js';
import './Game.scss';

export class Game extends Component {
	
	constructor(props) {
		super(props);
		
		this.play = this.play.bind(this);
	}

	play(column) {
		const { gameId, insertToken } = this.props;
		insertToken(gameId, column);
	}

	render() {
		const { board, cols, rows, newGame, k, level } = this.props;
		return (
			<div className='game'>
				<ConnectKBoard board={board} cols={cols} rows={rows} play={this.play} />
				<PropertyPanel newGame={newGame} cols={cols} rows={rows} k={k} level={level} />
			</div>
		);
	}
}

function mapStateToProps(state, ownProps) {
	return {
		gameId: state.board.gameId,
		board: state.board.board,
		cols: state.board.columns,
		rows: state.board.rows,
		k: state.board.k,
		level: state.board.level
	}
}


export default connect(mapStateToProps, {
	newGame,
	insertToken
})(Game)