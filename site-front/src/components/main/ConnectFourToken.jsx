import React, { Component } from 'react';
import classNames from 'classnames';

import './ConnectFourToken.scss';

export default class ConnectFourToken extends Component {

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
		const player = this.props.player;
		if (typeof player === 'undefined') {
			return null;
		}
		const classes = classNames('connectFourToken', { red: player == 'player', blue: player == 'computer' } );
		return (
			<div className={classes} />
		);
	}
}

ConnectFourToken.propTypes = {
	player: React.PropTypes.string,
}