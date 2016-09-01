import React, { Component } from 'react';
import { Provider } from 'react-redux';

import App from '../components/App.jsx';

export default class Root extends Component {
	
	render() {
		return (
			<Provider store={this.props.store}>
				<App />
			</Provider>
		);
	}
}

Root.propTypes = {
	store: React.PropTypes.object.isRequired
}