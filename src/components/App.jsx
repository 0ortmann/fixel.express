import React, { Component } from 'react';
import Header from './header/Header.jsx';
import Main from './main/Main.jsx';
import Footer from './footer/Footer.jsx';

import './App.scss';

const CONTENT = require('../content/content.json');

export default class App extends Component {
    
    constructor(props) {
        super(props);
    }

    render() {
        return (
            <div className='app'>
                <Header />
                <Main content={CONTENT} />
                <Footer />
            </div>
        );
    }
}