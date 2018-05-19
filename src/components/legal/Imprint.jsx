import React,{ Component } from 'react';
import { connect } from 'react-redux';
import Address from '../address/Address.jsx';
import Contact from '../address/Contact.jsx';

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
            <div className='content'>
                <h1>{imprint.title}</h1>
                {imprint.subtitle}
                <Address address={address} />
                <Contact contact={contact} />
                <h2>{imprint.sectionsTitle}</h2>
                {imprint.sections.map((section, index) => (
                    <div key={index} className='content__block'>
                        <h3>{section.title}</h3>
                        <p>{section.text}</p>
                    </div>
                ))};
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