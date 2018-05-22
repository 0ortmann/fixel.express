import React, { Component } from 'react';
import { connect } from 'react-redux';

import './AboutMe.scss';

export class AboutMe extends Component {
	
	constructor(props) {
		super(props);
	}

	render() {
		const { about } = this.props;
		return (
			<div>
				<div className='cover'>
					<div className='cover__title'>
						<h1>FELIX ORTMANN</h1>
						<p>{about.coverSubtitle}</p>
					</div>
				</div>
				<div className='profile container'>
					<h2>Profile</h2>
					<p>{about.profileSubtitle}</p>
					<hr />
					<div className='profile_content'>
						<div className='profile_content__description'>
							<h3>About Me</h3>
							<p>{about.description}</p>
						</div>
						<div className='profile_content__image'>
							<img src='/img/profile-small.png' alt='Felix Ortmann' width='320' height='320' />
						</div>
					</div>
				</div>
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