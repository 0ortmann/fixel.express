import React, { Component } from 'react';

import ConnectFourToken from './ConnectFourToken.jsx';

import './ConnectFourBoard.scss';

export default class ConnectFourBoard extends Component {

	constructor(props) {
		super(props);
		this.insert = this.insert.bind(this);
	}

	insert(evt) {
		const col = evt.target.dataset['col'];
		this.props.play(parseInt(col));
	}

	// TODO: dynamic column counts

	render() {
		const { board, cols, rows } = this.props;
		return (
			<div className='connectFourBoard'>
				{[...Array(cols)].map((x, c) => {
					return (
						<div data-col={c} key={c} className='connectFourBoard__column' onClick={this.insert}>
							{[...Array(rows)].map((y, r) => {
								return (
									<div data-col={c} key={r} className='connectFourBoard__row'>
										<ConnectFourToken player={board[c][rows - r -1]} />
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

ConnectFourBoard.propTypes = {
	cols: React.PropTypes.number.isRequired,
	rows: React.PropTypes.number.isRequired,
	play: React.PropTypes.func.isRequired,
	board: React.PropTypes.array.isRequired
}