import React, { Component } from 'react';
import { connect } from 'react-redux';

import './Imprint.scss';

export class Imprint extends Component {
    
    constructor(props) {
        super(props);
    }

    render() {
        const { imprint, address, contact } = this.props;
        if (!imprint) {
            return null;
        }
        return (
            <div className='imprint'>
                <h1>{imprint.title}</h1>
                {imprint.subtitle}
                <div className='imprint__block__fat'>
                    {address.name}<br />
                    {address.street}<br />
                    {address.zip} {address.city}, {address.country}
                </div>
                <div className='imprint__block__fat'>
                    {contact.phone}<br />
                    <a href={`mailto:${contact.email}`} target='_blank'>{contact.email}</a>
                </div>
                <h2>{imprint.sectionsTitle}</h2>
                {imprint.sections.map(section => {

                    return (<div className='imprint__block'>
                        <h3>{section.title}</h3>
                        <p>{section.text}</p>
                    </div>);
                })};
            </div>
        );
    }
}

function mapStateToProps(state) {
    return {
        imprint: state.lang.imprint,
        address: state.lang.address,
        contact: state.lang.contact,
    }
}

export default connect(mapStateToProps)(Imprint);