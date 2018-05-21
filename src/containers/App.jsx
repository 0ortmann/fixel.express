import React, { Component } from 'react';
import { connect } from 'react-redux';

import Footer from '../components/footer/Footer.jsx';

import { getLanguage } from '../actions/LanguageActionCreators.js';

import './App.scss';

const langs = ['german', 'english'];

export class App extends Component {

	static fetchData(store) {
		return store.dispatch(getLanguage('german'));
	}
	
	constructor(props) {
		super(props);
	}

	componentWillMount() {
		// fetch content
		this.props.getLanguage('german');
	}

	render() {
		const { getLanguage, children } = this.props;
		return (
			<div className='app'>
				{children}
				<Footer />
			</div>
		);
	}
}

function mapStateToProps(state, ownProps) {
	return {
		children: ownProps.children
	};
}

export default connect(mapStateToProps, {
	getLanguage
})(App);