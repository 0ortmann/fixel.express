import React, { Component } from 'react';

import './PropertyPanel.scss';

export default class PropertyPanel extends Component {
	
	constructor(props) {
		super(props);
		
		this.difficultyChange = this.difficultyChange.bind(this);
		this.columnChange = this.columnChange.bind(this);
		this.rowChange = this.rowChange.bind(this);
		this.newGame = this.newGame.bind(this);
	}

	componentWillMount() {
		this.props.newGame();
	}

	newGame(props) {
		const { newGame, cols, rows, k, level  } = this.props;
		let query = {
			c: cols,
			r: rows, 
			k: k,
			l: level
		};
		for (const prop in props) {
			query[prop] = props[prop];
		}
		newGame(query);
	}

	difficultyChange(evt) {
		evt.preventDefault();
		const level = evt.target.value;
		if (level) {
			this.newGame( { l: parseInt(level) } );
		}
	}

	columnChange(evt) {
		evt.preventDefault();
		const column = evt.target.value;
		if (column) {
			this.newGame( { c: parseInt(column) } );
		}
	}

	rowChange(evt) {
		evt.preventDefault();
		const row = evt.target.value;
		if (row) {
			this.newGame( { r: parseInt(row) } );
		}
	}

	render() {
		const { cols, rows, level, k } = this.props;
		return (
			<div className='game__properties'>
				<button className='game__new' onClick={this.newGame}>Ding Dong</button>
				<input className='range game__difficulty' type='range' min={2} max={8} step={2} defaultValue={level} onChange={this.difficultyChange} />
				<input className='number' type='number' min={7} max={20} step={1} defaultValue={cols} onChange={this.columnChange} />
				<input className='number' type='number' min={6} max={20} step={1} defaultValue={rows} onChange={this.rowChange} />
			</div>
		);
	}
}

PropertyPanel.propTypes = {
	newGame: React.PropTypes.func.isRequired,
	cols: React.PropTypes.number.isRequired,
	rows: React.PropTypes.number.isRequired,
	level: React.PropTypes.number.isRequired,
	k: React.PropTypes.number.isRequired
}