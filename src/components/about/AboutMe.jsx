import React, { Component } from 'react';
import { connect } from 'react-redux';

import './AboutMe.scss';

export class AboutMe extends Component {
	
	constructor(props) {
		super(props);
		this.state = { 
			hideFixHint: false
		};
	}

	render() {
		const { hideFixHint } = this.state;
		const { about } = this.props;
		return (
			<div className='about'>
				<h1>FELIX ORTMANN</h1>
				<div className='about__subtitle'>
					&mdash; {about.title} &mdash;
				</div>
				<img className='about__img'/>
				<div className='about__description'>
					{about.description.map((paragraph, i) => {
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
		about: state.lang.about
	}
}

export default connect(mapStateToProps)(AboutMe);