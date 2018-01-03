import React, { Component } from 'react';
import { connect } from 'react-redux';

import Footer from '../components/footer/Footer.jsx';
import Header from '../components/header/Header.jsx';
import LanguageSelector from '../components/language/LanguageSelector.jsx';

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
		const { langProps, getLanguage, children } = this.props;
		if (langProps == undefined || langProps == null) {
			return null;
		}
		return (
			<div className='app'>
				<Header boxText={langProps.box}/>
				<LanguageSelector selectCallback={getLanguage} langs={langs}/>
				{children}
				<Footer />
			</div>
		);
	}
}

function mapStateToProps(state, ownProps) {
	return {
		langProps: state.lang.properties,
		children: ownProps.children
	};
}

export default connect(mapStateToProps, {
	getLanguage
})(App);