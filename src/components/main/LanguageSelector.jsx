import React, {Component} from 'react';

import './LanguageSelector.scss';

export default class LanguageSelector extends Component {

	render() {
		const {langs, selectCallback} = this.props;
		return (
			<div className='languageSelector'>
				<select defaultValue={'german'} onClick={selectCallback}>
					{langs.map(lang => {
						return <option key={lang}>{lang}</option>
					})}
				</select>
			</div>
		);
	}

}

LanguageSelector.proptypes = {
	selectCallback: React.PropTypes.func.isRequired,
	langs: React.PropTypes.array.isRequired
}