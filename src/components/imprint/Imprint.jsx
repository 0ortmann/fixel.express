import React, { Component } from 'react';
import { connect } from 'react-redux';

import './Imprint.scss';

export class Imprint extends Component {
    
    constructor(props) {
        super(props);
    }

    render() {
        const { imprint } = this.props;
        if (!imprint) {
            return null;
        }
        return (
            <div className='imprint'>
                <h1>{imprint.title}</h1>
                {imprint.subtitle}
                <div className='imprint__block__fat'>
                    {imprint.address.name}<br />
                    {imprint.address.street}<br />
                    {imprint.address.zip} {imprint.address.city}, {imprint.address.country}
                </div>
                <div className='imprint__block__fat'>
                    {imprint.contact.phone}<br />
                    <a href={`mailto:${imprint.contact.email}`} target='_blank'>{imprint.contact.email}</a>
                </div>
                <h2>{imprint.disclaimerTitle}</h2>
                <div className='imprint__block'>
                    <h3>{imprint.contentAccountability.title}</h3>
                    <p>{imprint.contentAccountability.text}</p>
                </div>
                <div className='imprint__block'>
                    <h3>{imprint.linkAccountability.title}</h3>
                    <p>{imprint.linkAccountability.text}</p>
                </div>
                <div className='imprint__block'>
                    <h3>{imprint.copyright.title}</h3>
                    <p>{imprint.copyright.text}</p>
                </div>
            </div>
        );
    }
}

function mapStateToProps(state) {
    return {
        imprint: state.lang.imprint
    }
}

export default connect(mapStateToProps)(Imprint);