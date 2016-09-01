import React from 'react';
import ReactDOM from 'react-dom';
import TestUtils from 'react-addons-test-utils';
import FixHeaderBox from '../src/components/main/FixHeaderBox.jsx';

describe('FixHeaderBox', () => {

	const text = 'DISPLAY TEXT';
	const cb = jest.genMockFunction();

	let box;

	beforeEach(function() {
        box = TestUtils.renderIntoDocument( <FixHeaderBox displayText={text} fixIt={cb} hide={false} /> );
    });

	it('should render', () => {
        let boxNode = ReactDOM.findDOMNode(box);
        expect(boxNode).toBeDefined();
        let renderedText = TestUtils.findRenderedDOMComponentWithTag(box, 'h3');
        expect(renderedText.textContent).toEqual(text);
	});

	it('should render correct hide-class', () => {
        box = TestUtils.renderIntoDocument( <FixHeaderBox displayText={text} fixIt={cb} hide={true} /> );
        let hideNode = TestUtils.findRenderedDOMComponentWithClass(box, 'hide');
        expect(hideNode).toBeDefined();
	});


	it('call the callback to fix header', () => {
		let link = TestUtils.findRenderedDOMComponentWithTag(box, 'a');
		expect(link).toBeDefined()
		TestUtils.Simulate.click( link );
		expect(cb).toBeCalled();
	});
});