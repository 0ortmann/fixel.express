import React, { Component } from 'react';
import './Header.scss';

export default class Header extends Component {
    
    constructor(props) {
        super(props);
    }

    render() {
        return (
            <header>
                <div className='content clearfix'>
                    <div className='logo'>
                        hier ein logo
                    </div>
                    <nav>
                        <ul>
                            <li>About</li>
                            <li>Projekte</li>
                            <li>Contact</li>
                        </ul>
                    </nav>
                    <div className='social'>
                        hier social links
                    </div>
                </div>
            </header>
        );
    }
}