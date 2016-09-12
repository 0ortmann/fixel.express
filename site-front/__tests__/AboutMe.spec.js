import React from 'react';
import ReactDOM from 'react-dom';
import TestUtils from 'react-addons-test-utils';

// dont import the connected ('wrapped') component, but the pure one
import { AboutMe } from '../src/components/about/AboutMe.jsx';

describe('AboutMe', () => {

	const gerCont = { 
		description: ['foo', 'bar'],
		title: 't1',
		box: 'box'
	}; 
	const engCont = {
		description: ['ping', 'pong'],
		title: 't2',
		box: 'box'
	};

	const fixHeader = jest.genMockFunction();
	const getLanguage = jest.genMockFunction();

	let about;


	beforeEach(function() {
		about = TestUtils.renderIntoDocument(
			<AboutMe langProps={gerCont} getLanguage={getLanguage} fixHeader={fixHeader} />
		);
	});

	it('should render', () => {
		const aboutNode = ReactDOM.findDOMNode(about);
		expect(aboutNode).toBeDefined();
		
		const paragraphs = TestUtils.scryRenderedDOMComponentsWithTag(about, 'p');
		expect(paragraphs.length).toEqual(2);
		expect(paragraphs[0].textContent).toEqual('foo');
		expect(paragraphs[1].textContent).toEqual('bar');
	});
});