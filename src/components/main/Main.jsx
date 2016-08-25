import React, { Component } from 'react';

import LanguageSelector from './LanguageSelector.jsx';
import './Main.scss';

const langs = ['german', 'english']

export default class Main extends Component {
    
    constructor(props) {
        super(props);
        this.switchLanguage = this.switchLanguage.bind(this);
        this.state = { content: undefined };
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

    render() {
        if (this.state.content == undefined) {
            return null;
        }
        return (
            <main className='about'>
                <LanguageSelector selectCallback={this.switchLanguage} langs={langs}/>
                <h1>FELIX ORTMANN</h1>
                <div className='about__subtitle'>
                    &mdash; {this.state.content.title} &mdash;
                </div>
                <img className='about__img'/>
                <div className='about__description'>
                    {this.state.content.description.map((paragraph, i) => {
                        return <p key={i}>{paragraph}</p>
                    })}
                </div>
            </main>
        );
    }
}

Main.propTypes = {
    content: React.PropTypes.object.isRequired
}