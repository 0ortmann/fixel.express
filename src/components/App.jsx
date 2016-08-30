import React, { Component } from 'react';
import Ugly from './header/Ugly.jsx';
import Hammer from './header/Hammer.jsx';
import Shiny from './header/Shiny.jsx';
import Main from './main/Main.jsx';
import Footer from './footer/Footer.jsx';

import './App.scss';

const CONTENT = require('../content/content.json');

export default class App extends Component {
    
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
            <div className='app'>
                {this.getHeader()}
                <Main content={CONTENT} fixHeader={this.hammerHeader} />
                <Footer />
            </div>
        );
    }
}