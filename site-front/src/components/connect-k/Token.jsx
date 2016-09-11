import React, { Component } from 'react';
import classNames from 'classnames';

import './Token.scss';

export default class Token extends Component {

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
		const { player, col } = this.props;
		if (typeof player === 'undefined') {
			return null;
		}
		const classes = classNames('token', { red: player == 'player', blue: player == 'computer' } );
		return (
			<div className={classes} data-col={col}/>
		);
	}
}

Token.propTypes = {
	player: React.PropTypes.string,
	col: React.PropTypes.number
}