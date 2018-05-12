import React, {Component} from 'react';
import PropTypes from 'prop-types';
import classNames from 'classnames';

import './FixHeaderBox.scss';


export default class FixHeaderBox extends Component {

	render() {
		const { fixIt, hide, displayText } = this.props;
		const myClass = classNames('fixHeaderBox', {hide: hide});
		return (
			<div className={myClass}>
				<object className='fix__arrow'/>
				<h3>{displayText}</h3>
				<a onClick={fixIt}>Fix it!</a>
			</div>
		);
	}
}

FixHeaderBox.propTypes = {
	displayText: PropTypes.string.isRequired,
	fixIt: PropTypes.func.isRequired,
	hide: PropTypes.bool.isRequired
}