import React, { Component } from 'react';

import './PropertyPanel.scss';

export default class PropertyPanel extends Component {
	
	constructor(props) {
		super(props);
		
		this.difficultyChange = this.difficultyChange.bind(this);
		this.columnChange = this.columnChange.bind(this);
		this.rowChange = this.rowChange.bind(this);
		this.kChange = this.kChange.bind(this);
		this.newButton = this.newButton.bind(this);
		
		this.newGame = this.newGame.bind(this);
	}

	componentWillMount() {
		this.props.newGame();
	}

	newButton() {
		// need to swallow the event...
		this.newGame({});
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
			if(props.hasOwnProperty(prop)) {
				query[prop] = props[prop];
			}
		}
		newGame(query);
	}

	handleChange(evt, key) {
		evt.preventDefault();
		const val = evt.target.value;
		if (val) {
			let q = {};
			q[key] = parseInt(val);
			this.newGame(q);
		}
	}

	difficultyChange(evt) {
		this.handleChange(evt, 'l');
	}

	columnChange(evt) {
		this.handleChange(evt, 'c');
	}

	rowChange(evt) {
		this.handleChange(evt, 'r');
	}

	kChange(evt) {
		this.handleChange(evt, 'k');
	}

	render() {
		const { cols, rows, level, k } = this.props;
		return (
			<div className='game__properties'>
				<button className='game__new' onClick={this.newButton}>Ding Dong</button>
				<input className='range game__difficulty' type='range' min={2} max={8} step={2} defaultValue={level} onChange={this.difficultyChange} />
				<input className='number' type='number' min={7} max={20} step={1} defaultValue={cols} onChange={this.columnChange} />
				<input className='number' type='number' min={6} max={20} step={1} defaultValue={rows} onChange={this.rowChange} />
				<input className='number' type='number' min={3} max={8} step={1} defaultValue={k} onChange={this.kChange} />
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