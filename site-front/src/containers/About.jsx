import React, { Component } from 'react';
import Ugly from '../components/header/Ugly.jsx';
import Hammer from '../components/header/Hammer.jsx';
import Shiny from '../components/header/Shiny.jsx';
import AboutMe from '../components/about/AboutMe.jsx';

export default class About extends Component {
	
	constructor(props) {
		super(props);
		this.hammerHeader = this.hammerHeader.bind(this);
		this.shinyHeader = this.shinyHeader.bind(this);

		this.state = { header: 'ugly' };
	}

	hammerHeader() {
		this.setState({ header: 'hammer' });
	}

	shinyHeader() {
		console.log("shiny")
		this.setState({ header: 'shiny' });
	}

	getHeader() {
		if (this.state.header == 'ugly') {
			return <Ugly />;
		}
		else if (this.state.header == 'hammer') {
			return <Hammer doneHammering={this.shinyHeader} />;
		}
		return <Shiny />;
	}

	render() {
		return (
			<div className='about'>
				{this.getHeader()}
				<AboutMe fixHeader={this.hammerHeader} />
			</div>
		);
	}
}