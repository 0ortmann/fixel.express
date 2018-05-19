import React,{ Component } from 'react';
import PropTypes from 'prop-types';

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
            <div className='content__block'>
                <h3>{title}</h3>
                <p>{text.map(subsection => (<span>{subsection}<br /></span>))}</p>
            </div>
        );
    }
}

LegalSection.propTypes = {
    title: PropTypes.string.isRequired,
    text: PropTypes.array.isRequired
}