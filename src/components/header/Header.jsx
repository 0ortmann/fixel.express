import React, { Component } from 'react';
import './Header.scss';

export default class Header extends Component {
    
    constructor(props) {
        super(props);
    }

    render() {
        return (
            <header>
                <div className='header__fixel'>
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
            </header>
        );
    }
}