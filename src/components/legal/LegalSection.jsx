import React,{ Component } from 'react';
import PropTypes from 'prop-types';

import './Legal.scss';

export default class LegalSection extends Component {
    
    constructor(props) {
        super(props);
    }

    render() {
        const { title, text } = this.props;
        if (!title || !text) {
            return null;
        }
        return (
            <div className='legal_content__block'>
                <h4>{title}</h4>
                <p>{text.map(subsection => (<span>{subsection}<br /></span>))}</p>
            </div>
        );
    }
}

LegalSection.propTypes = {
    title: PropTypes.string.isRequired,
    text: PropTypes.array.isRequired
}