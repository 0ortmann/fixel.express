import React, { Component } from 'react';
import classNames from 'classnames';

import './Shiny.scss';

export default class Shiny extends Component {
    
    constructor(props) {
        super(props);
    }


    render() {
        return (
            <header>
                <h1 className='shinyHeader'>I'm not a design person</h1>
            </header>
        );
    }
}
