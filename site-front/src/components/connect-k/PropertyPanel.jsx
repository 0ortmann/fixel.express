import React, { Component } from 'react';

import './PropertyPanel.scss';

export default class PropertyPanel extends Component {
	
	constructor(props) {
		super(props);
		
		this.difficultyChange = this.difficultyChange.bind(this);
		this.newGame = this.newGame.bind(this);
	}

	componentWillMount() {
		this.props.newGame();
	}

	newGame() {
		const { newGame, cols, rows, k, level  } = this.props;
		newGame({
			c: cols,
			r: rows, 
			k: k,
			l: level
		});

	}

	difficultyChange(evt) {
		const level = evt.target.value;
		if (level) {
			this.props.newGame( { l: parseInt(level) } );
		}
	}

	render() {
		const { cols, rows, level, k } = this.props;
		return (
			<div className='game__properties'>
				<button className='game__new' onClick={this.newGame}>Ding Dong</button>
				<input className='range game__difficulty' type='range' min={2} max={8} step={2} defaultValue={level} onChange={this.difficultyChange} />
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