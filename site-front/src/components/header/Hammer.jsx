import React, { Component } from 'react';
import classNames from 'classnames';

import './Hammer.scss';

export default class Hammer extends Component {
    
    constructor(props) {
        super(props);
        
        this.hammer = this.hammer.bind(this);
        this.state = { hammerUp : false };

        // has to match css... not that nice.
        const wait = 1500;
        const hammerTime = 2500;
        
        setTimeout(() => {
            const hammerInterval = setInterval(this.hammer, 100);
            setTimeout(() => clearInterval(hammerInterval), hammerTime);
        }, wait);

        setTimeout(() => {
            this.props.doneHammering();
        }, wait + hammerTime)
    }

    hammer() {
        this.setState({hammerUp: !this.state.hammerUp})
    }

    render() {
        const gopherLeft = classNames('gopherHammer', { 'gopherHammer-left': this.state.hammerUp, 'gopherHammer-left2': !this.state.hammerUp});
        const gopherRight = classNames('gopherHammer', { 'gopherHammer-right': this.state.hammerUp, 'gopherHammer-right2': !this.state.hammerUp});
        return (
            <header>
                <div className='gopherContainer move-right'>
                    <object className={gopherLeft} />
                </div>
                <div className='expressWrapper'>
                    <div className='fixel__smoke'> 
                        <span className='fixel__smoke__0'></span>
                        <span className='fixel__smoke__1'></span>
                        <span className='fixel__smoke__2'></span>
                        <span className='fixel__smoke__3'></span>
                    </div>
                    <div className='fixel__jumping'>F</div>
                    <div className='fixel__jumping'>I</div>
                    <div className='fixel__jumping'>X</div>
                    <div className='fixel__jumping'>E</div>
                    <div className='fixel__jumping'>L</div>

                    <div className='express__e1 fixel__jumping'>E</div>
                    <div className='fixel__jumping'>X</div>
                    <div className='fixel__jumping'>P</div>
                    <div className='fixel__jumping'>R</div>
                    <div className='fixel__jumping'>E</div>
                    <div className='fixel__jumping'>S</div>
                    <div className='fixel__jumping'>S</div>
                </div>
                <div className='gopherContainer move-left'>
                    <object className={gopherRight} />
                </div>
            </header>
        );
    }
}

Hammer.propTypes = {
    doneHammering: React.PropTypes.func.isRequired
}