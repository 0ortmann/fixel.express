import React, { Component } from 'react';
import RangeInput from './RangeInput.jsx';

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
		const { newGame, cols, rows, k, level } = this.props;
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

	handleChange(val, key) {
		if (val) {
			let q = {};
			q[key] = parseInt(val);
			this.newGame(q);
		}
	}

	difficultyChange(val) {
		this.handleChange(val, 'l');
	}

	columnChange(val) {
		this.handleChange(val, 'c');
	}

	rowChange(val) {
		this.handleChange(val, 'r');
	}

	kChange(val) {
		this.handleChange(val, 'k');
	}

	render() {
		const { cols, rows, level, k, propertyLabels } = this.props;
		return (
			<div className='game__properties'>
				<button className='game__new' onClick={this.newButton}>{propertyLabels.new}</button>
				<RangeInput step={2} min={2} max={8} callback={this.difficultyChange} defaultValue={level} label={propertyLabels.level}/>
				<RangeInput step={1} min={3} max={11} callback={this.columnChange} defaultValue={cols} label={propertyLabels.cols}/>
				<RangeInput step={1} min={3} max={11} callback={this.rowChange} defaultValue={rows} label={propertyLabels.rows}/>
				<RangeInput step={1} min={3} max={8} callback={this.kChange} defaultValue={k} label={propertyLabels.k}/>
			</div>
		);
	}
}

PropertyPanel.propTypes = {
	newGame: React.PropTypes.func.isRequired,
	cols: React.PropTypes.number.isRequired,
	rows: React.PropTypes.number.isRequired,
	level: React.PropTypes.number.isRequired,
	k: React.PropTypes.number.isRequired,
	propertyLabels: React.PropTypes.object
}