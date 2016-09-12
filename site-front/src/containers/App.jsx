import React, { Component } from 'react';
import { connect } from 'react-redux';

import Footer from '../components/footer/Footer.jsx';
import Header from '../components/header/Header.jsx';
import { getLanguage } from '../actions/LanguageActionCreators.js';

import './App.scss';


export class App extends Component {
	
	constructor(props) {
		super(props);
	}

	componentWillMount() {
		this.props.getLanguage('german');
	}

	render() {
		return (
			<div className='app'>
				<Header />
				{this.props.children}
				<Footer />
			</div>
		);
	}
}

export default connect(null, {
	getLanguage
})(App);