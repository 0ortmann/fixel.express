import React from 'react';
import ReactDOM from 'react-dom';
import TestUtils from 'react-addons-test-utils';
import Hammer from '../../src/components/header/Hammer.jsx';

describe('Hammer', () => {

	const doneCb = jest.genMockFunction();

	let hammer;

	beforeEach(function() {
		hammer = TestUtils.renderIntoDocument( <Hammer doneHammering={doneCb} /> );
	});

	it('should render', () => {
		let hammerNode = ReactDOM.findDOMNode(hammer);
		expect(hammerNode).toBeDefined();

		expect(hammer.state.hammerUp).toEqual(false);
		const gopherLeft = TestUtils.findRenderedDOMComponentWithClass(hammer, 'gopherHammer-left2');
		const gopherRight = TestUtils.findRenderedDOMComponentWithClass(hammer, 'gopherHammer-right2');

		expect(gopherLeft).toBeDefined();
		expect(gopherRight).toBeDefined();
	});

	it('should change hammer state', () => {
		hammer.hammer();

		expect(hammer.state.hammerUp).toEqual(true);

		const gopherLeft = TestUtils.findRenderedDOMComponentWithClass(hammer, 'gopherHammer-left');
		const gopherRight = TestUtils.findRenderedDOMComponentWithClass(hammer, 'gopherHammer-right');

		expect(gopherLeft).toBeDefined();
		expect(gopherRight).toBeDefined();
	});
});