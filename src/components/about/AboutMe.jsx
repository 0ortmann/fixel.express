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
		const { langProps } = this.props;
		return (
			<div className='about'>
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

export default connect(mapStateToProps)(AboutMe);