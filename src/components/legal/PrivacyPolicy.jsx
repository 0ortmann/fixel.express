import React,{ Component } from 'react';
import { connect } from 'react-redux';
import Address from '../address/Address.jsx';
import Contact from '../address/Contact.jsx';

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
            <div className='content'>
                <h1>{privPolicy.title}</h1>
                {privPolicy.subtitle}
                <Address address={address} />
                <Contact contact={contact} />
                <h2>{privPolicy.sectionsTitle}</h2>
                {privPolicy.sections.map((section, index) => (
                    <div key={index} className='content__block'>
                        <h3>{section.title}</h3>
                        <p>{section.text.map(subsection => (<span>{subsection}<br /></span>))}</p>
                    </div>
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