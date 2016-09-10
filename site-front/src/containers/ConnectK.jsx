import React, { Component } from 'react';

import Game from '../components/connect-k/Game.jsx';

export default class ConnectK extends Component {
	
	constructor(props) {
		super(props);
	}

	render() {
		return (
			<div className='connectK'>
				<Game />
			</div>
		);
	}
}