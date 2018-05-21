import React,{ Component } from 'react';
import { connect } from 'react-redux';
import Address from '../address/Address.jsx';
import Contact from '../address/Contact.jsx';
import LegalSection from './LegalSection.jsx';

import './Legal.scss';

export class PrivacyPolicy extends Component {
    
    constructor(props) {
        super(props);
    }

    render() {
        const { privPolicy, address, contact } = this.props;
        if (!privPolicy) {
            return null;
        }
        return (
            <div className='legal_container'>
                <h2>{privPolicy.title}</h2>
                <h4>{privPolicy.subtitle}</h4>
                <Address address={address} />
                <Contact contact={contact} />
                <h3>{privPolicy.sectionsTitle}</h3>
                {privPolicy.sections.map((section, index) => (
                    <LegalSection key={index} title={section.title} text={section.text} />
                ))};
            </div>
        );
    }
}

function mapStateToProps(state) {
    return {
        privPolicy: state.lang.privPolicy,
        address: state.lang.address,
        contact: state.lang.contact,
    }
}

export default connect(mapStateToProps)(PrivacyPolicy);