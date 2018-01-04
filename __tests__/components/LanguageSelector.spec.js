import React from 'react';
import ReactDOM from 'react-dom';
import TestUtils from 'react-dom/test-utils';
import LanguageSelector from '../../src/components/language/LanguageSelector.jsx';

describe('LanguageSelector', () => {

	const langs = ['l1', 'l2', 'l3'];
	const cb = jest.genMockFunction();

	let selector;

	beforeEach(function() {
		selector = TestUtils.renderIntoDocument( <LanguageSelector langs={langs} selectCallback={cb} /> );
	});

	it('should render', () => {
		let selectorNode = ReactDOM.findDOMNode(selector);
		expect(selectorNode).toBeDefined();
		let options = TestUtils.scryRenderedDOMComponentsWithTag(selector, 'option');
		expect(options.length).toEqual(3);
	});

	it('call the callback with language paramter', () => {
		let select = TestUtils.findRenderedDOMComponentWithTag(selector, 'select');
		expect(select).toBeDefined();
		TestUtils.Simulate.change( select, {target: {value: 'l2' }} );
		expect(cb).toBeCalledWith('l2');
	});

});