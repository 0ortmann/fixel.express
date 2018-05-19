import React, { Component } from 'react';
import PropTypes from 'prop-types';

export default class Contact extends Component {
    
    constructor(props) {
        super(props);
    }

    render() {
        const { contact } = this.props;
        if (!contact) {
            return null;
        }
        return (
            <div className='content__block__fat'>
                {contact.phone}<br />
                <a href={`mailto:${contact.email}`} target='_blank'>{contact.email}</a>
            </div>
        );
    }
}
Contact.propTypes = {
    contact: PropTypes.object.isRequired
}

