import React, { Component } from 'react';
import { connect } from 'react-redux';

import ConnectKBoard from './ConnectKBoard.jsx';
import PropertyPanel from './PropertyPanel.jsx';

import { newGame, insertToken } from '../../actions/ConnectKActionCreators.js';
import './Game.scss';

export class Game extends Component {
	
	constructor(props) {
		super(props);
		
		this.play = this.play.bind(this);
	}

	play(column) {
		const { id, insertToken } = this.props;
		insertToken(id, column);
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
		id: state.game.id,
		board: state.game.board,
		cols: state.game.columns,
		rows: state.game.rows,
		k: state.game.k,
		level: state.game.level
	}
}


export default connect(mapStateToProps, {
	newGame,
	insertToken
})(Game)