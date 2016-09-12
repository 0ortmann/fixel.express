import React from 'react';
import ReactDOM from 'react-dom';
import TestUtils from 'react-addons-test-utils';

import { Header } from '../src/components/header/Header.jsx';

describe('Header', () => {

	const langProps = { box: 'BOX TEXT'};
	let header;

	beforeEach(function() {
		header = TestUtils.renderIntoDocument( <Header langProps={langProps} /> );
	});

	it('should render with ugly header as default', () => {
		const headerNode = ReactDOM.findDOMNode(header);
		expect(headerNode).toBeDefined();
		
		expect(header.state.header).toEqual('ugly');
		expect(header.state.hideFixHint).toEqual(false);
	});
});