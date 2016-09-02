import React, { Component } from 'react';

//import './ConnectFourToken.scss';

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
		if (typeof player === 'undefiend') {
			return null;
		}
		return (
			<div className='connectFourToken'>
				{player}
			</div>

		);
	}
}

ConnectFourToken.propTypes = {
	player: React.PropTypes.string,
}