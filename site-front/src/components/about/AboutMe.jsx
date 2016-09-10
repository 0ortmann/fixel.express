import React, { Component } from 'react';

import LanguageSelector from './LanguageSelector.jsx';
import FixHeaderBox from './FixHeaderBox.jsx';

import './AboutMe.scss';

const langs = ['german', 'english']

export default class AboutMe extends Component {
	
	constructor(props) {
		super(props);
		this.state = { 
			content: undefined,
			hideFixHint: false
		};
		this.switchLanguage = this.switchLanguage.bind(this);
		this.fixHeader = this.fixHeader.bind(this);
	}

	componentWillMount() {
		this.switchLanguage('german');
	}

	switchLanguage(lang) {
		if (lang == undefined || lang == '') {
			return;
		}
		const use = this.props.content[lang]
		if (use != undefined) {
			this.setState( {'content': use});
		}
	}

	fixHeader() {
		this.setState({ hideFixHint: true });
		this.props.fixHeader();
	}

	render() {
		const { content, hideFixHint } = this.state;
		const { board, cols, rows, newGame } = this.props;
		if (content == undefined) {
			return null;
		}
		return (
			<div className='about'>
				<LanguageSelector selectCallback={this.switchLanguage} langs={langs}/>
				<h1>FELIX ORTMANN</h1>
				<div className='about__subtitle'>
					&mdash; {content.title} &mdash;
				</div>
				<FixHeaderBox fixIt={this.fixHeader} hide={hideFixHint} displayText={content.box} />
				<img className='about__img'/>
				<div className='about__description'>
					{content.description.map((paragraph, i) => {
						return <p key={i}>{paragraph}</p>
					})}
				</div>
				<object className='gopher__left'/>
			</div>
		);
	}
}

AboutMe.propTypes = {
	content: React.PropTypes.object.isRequired,
	fixHeader: React.PropTypes.func.isRequired
}