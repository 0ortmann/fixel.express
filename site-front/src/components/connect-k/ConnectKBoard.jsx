import React, { Component } from 'react';

import Token from './Token.jsx';

import './ConnectKBoard.scss';

export default class ConnectKBoard extends Component {

	constructor(props) {
		super(props);
		this.insert = this.insert.bind(this);
	}

	insert(evt) {
		const col = evt.target.dataset['col'];
		if (col) {
			this.props.play(parseInt(col));
		}
	}

	// TODO: dynamic column counts

	render() {
		const { board, cols, rows } = this.props;
		const boardStyle = { columnCount: cols, 'maxWidth': cols * 85, 'minWidth': cols * 75 };
		return (
			<div className='connectKBoard' style={boardStyle}>
				{[...Array(cols)].map((x, c) => {
					return (
						<div data-col={c} key={c} className='connectKBoard__column'>
							{[...Array(rows)].map((y, r) => {
								return (
									<div data-col={c} key={r} className='connectKBoard__row' onClick={this.insert}>
										<Token player={board[c][rows - r -1]} col={c} />
									</div>
								);
							})}
						</div>
					);
				})}
			</div>

		);
	}
}

ConnectKBoard.propTypes = {
	cols: React.PropTypes.number.isRequired,
	rows: React.PropTypes.number.isRequired,
	play: React.PropTypes.func.isRequired,
	board: React.PropTypes.array.isRequired
}