import React, { Component } from 'react';
import { connect } from 'react-redux';

import { getLanguage } from '../../actions/LanguageActionCreators.js';

import LanguageSelector from './LanguageSelector.jsx';

import './AboutMe.scss';

const langs = ['german', 'english']

export class AboutMe extends Component {
	
	constructor(props) {
		super(props);
		this.state = { 
			hideFixHint: false
		};
	}

	render() {
		const { hideFixHint } = this.state;
		const { langProps } = this.props;
		if (langProps == undefined) {
			return null;
		}
		return (
			<div className='about'>
				<LanguageSelector selectCallback={this.props.getLanguage} langs={langs}/>
				<h1>FELIX ORTMANN</h1>
				<div className='about__subtitle'>
					&mdash; {langProps.title} &mdash;
				</div>
				<img className='about__img'/>
				<div className='about__description'>
					{langProps.description.map((paragraph, i) => {
						return <p key={i}>{paragraph}</p>
					})}
				</div>
				<object className='gopher__left'/>
			</div>
		);
	}
}

function mapStateToProps(state) {
	return {
		langProps: state.lang.properties
	}
}

export default connect(mapStateToProps, {
	getLanguage
})(AboutMe);