import React, {Component} from 'react';

import './FixHeaderBox.scss';

export default class FixHeaderBox extends Component {

	render() {
		return (
			<div className='fixHeaderBox'>
				<object className='fix__arrow'/>
				<h3>Sieht voll kacke aus?</h3>
				<a href="">Fix it!</a>
			</div>
		);
	}
}