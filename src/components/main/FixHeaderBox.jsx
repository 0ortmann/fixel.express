import React, {Component} from 'react';
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
	displayText: React.PropTypes.string.isRequired,
	fixIt: React.PropTypes.func.isRequired,
	hide: React.PropTypes.bool.isRequired
}