import React, { Component } from 'react';

import './RangeInput.scss';

export default class RangeInput extends Component {

	constructor(props) {
		super(props);
		this.increase = this.increase.bind(this);
		this.decrease = this.decrease.bind(this);
	}

	decrease(evt) {
		const { callback, defaultValue, step, min } = this.props;
		const val = defaultValue >= min + step? defaultValue - step : defaultValue;
		callback(val);
	}

	increase(evt) {
		const { callback, defaultValue, step, max } = this.props;
		const val = defaultValue <= max - step ? defaultValue + step : defaultValue;
		console.log(val)
		callback(val);
	}

	render() {
		const {step, min, max, defaultValue, callback, label} = this.props;
		let htmlLabel = '';
		if (label) {
			htmlLabel = <label>{label}:</label>;
		}
		return (
			<div>
				{htmlLabel}
				<button className='rangeInput__arrow' type='button' onClick={this.decrease}>&lsaquo;</button>
				<span className='rangeInput__field'>{defaultValue} </span>
				<button className='rangeInput__arrow' type='button' onClick={this.increase}>&rsaquo;</button>
			</div>
		);
	}
}

RangeInput.propTypes = {
	step: React.PropTypes.number.isRequired,
	min: React.PropTypes.number.isRequired,
	max: React.PropTypes.number.isRequired,
	defaultValue: React.PropTypes.number.isRequired,
	callback: React.PropTypes.func.isRequired,
	label: React.PropTypes.string
}