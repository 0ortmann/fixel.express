import React from 'react';
import ReactDOM from 'react-dom';
import TestUtils from 'react-dom/test-utils';

import Header from '../../src/components/header/Header.jsx';

describe('Header', () => {

	const boxText = 'BOX TEXT';
	let header;

	beforeEach(function() {
		header = TestUtils.renderIntoDocument( <Header boxText={boxText} /> );
	});

	it('should render with ugly header as default', () => {
		const headerNode = ReactDOM.findDOMNode(header);
		expect(headerNode).toBeDefined();
		
		expect(header.state.header).toEqual('ugly');
		expect(header.state.hideFixHint).toEqual(false);
	});
});