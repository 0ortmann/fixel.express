import React, { Component } from 'react';
import PropTypes from 'prop-types';

export default class Address extends Component {
    
    constructor(props) {
        super(props);
    }

    render() {
        const { address } = this.props;
        if (!address) {
            return null;
        }
        return (
            <div className='content__block__fat'>
                {address.name}<br />
                {address.street}<br />
                {address.zip} {address.city}, {address.country}
            </div>
        );
    }
}
Address.propTypes = {
    address: PropTypes.object.isRequired
}

