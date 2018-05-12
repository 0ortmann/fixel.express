import PropTypes from 'prop-types';
import React, { Component } from 'react';

import Ugly from './Ugly.jsx';
import Hammer from './Hammer.jsx';
import Shiny from './Shiny.jsx';

import FixHeaderBox from './FixHeaderBox.jsx';

import './Header.scss';

// Kinda the wrapper for all the header fun!
export default class Header extends Component {
	
	constructor(props) {
		super(props);
		this.hammerHeader = this.hammerHeader.bind(this);
		this.shinyHeader = this.shinyHeader.bind(this);

		this.state = { 
			header: 'ugly',
			hideFixHint: false
		};
	}

	hammerHeader() {
		this.setState({ hideFixHint: true });
		this.setState({ header: 'hammer' });
	}

	shinyHeader() {
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
		const { boxText } = this.props;
		const hideFixHint = this.state.hideFixHint;
		return (
			<div>
			{this.getHeader()}
			<FixHeaderBox displayText={boxText} fixIt={this.hammerHeader} hide={this.state.hideFixHint} />
			</div>
		);
	}
}

Header.propTypes = {
	boxText: PropTypes.string.isRequired
}