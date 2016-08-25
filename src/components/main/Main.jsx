import React, { Component } from 'react';

import LanguageSelector from './LanguageSelector.jsx';
import './Main.scss';

const CONTENT = require('../../content/content.json');
const langs = ['german', 'english']

export default class Main extends Component {
    
    constructor(props) {
        super(props);
        this.state = { content: CONTENT['german'] };
        this.switchLanguage = this.switchLanguage.bind(this);
    }

    switchLanguage(evt) {
        const lang = evt.target.value;
        if (lang == undefined) {
            return;
        }
        const use = CONTENT[lang]
        if (use != undefined) {
            this.setState( {'content': use});
        }
    }

    render() {
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