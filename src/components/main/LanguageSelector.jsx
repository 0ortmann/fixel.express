import React, {Component} from 'react';

import './LanguageSelector.scss';

export default class LanguageSelector extends Component {

	handleSelect(evt) {
		evt.preventDefault();
		this.props.selectCallback(evt.target.value);
	}

	render() {
		const {langs, selectCallback} = this.props;
		return (
			<div className='languageSelector'>
				<select defaultValue={langs[0]} onChange={this.handleSelect.bind(this)}>
					{langs.map(lang => {
						return <option key={lang}>{lang}</option>
					})}
				</select>
			</div>
		);
	}

}

LanguageSelector.propTypes = {
	selectCallback: React.PropTypes.func.isRequired,
	langs: React.PropTypes.array.isRequired
}