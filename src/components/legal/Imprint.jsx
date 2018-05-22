import React,{ Component } from 'react';
import { connect } from 'react-redux';
import Address from '../address/Address.jsx';
import Contact from '../address/Contact.jsx';
import LegalSection from './LegalSection.jsx';

import './Legal.scss';

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
            <div className='legal_container'>
                <h2>{imprint.title}</h2>
                <h4>{imprint.subtitle}</h4>
                <Address address={address} />
                <Contact contact={contact} />
                <h3>{imprint.sectionsTitle}</h3>
                {imprint.sections.map((section, index) => (
                    <LegalSection key={index} title={section.title} text={section.text} />
                ))}
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