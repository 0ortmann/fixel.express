import React, { Component } from 'react';
import Footer from '../components/footer/Footer.jsx';

import './App.scss';

export default class App extends Component {
	
	constructor(props) {
		super(props);
	}

	render() {
		return (
			<div className='app'>
				{this.props.children}
				<Footer />
			</div>
		);
	}
}