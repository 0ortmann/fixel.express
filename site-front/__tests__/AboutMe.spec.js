import React from 'react';
import ReactDOM from 'react-dom';
import TestUtils from 'react-addons-test-utils';

// dont import the connected ('wrapped') component, but the pure one
import AboutMe from '../src/components/about/AboutMe.jsx';

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

	const content = {
		german: gerCont,
		english: engCont
	};

	const fixHeader = jest.genMockFunction();
	const newGame = jest.genMockFunction();
	const cols = 7;
	const rows = 6;
	const gameId = '123';
	const board = [[], [], [], [], [], [], []];

	let about;


	beforeEach(function() {
		about = TestUtils.renderIntoDocument(
			<AboutMe content={content} fixHeader={fixHeader} />
		);
	});

	it('should render', () => {
		const aboutNode = ReactDOM.findDOMNode(about);
		expect(aboutNode).toBeDefined();
		expect(about.state.content).toEqual(gerCont);
		const paragraphs = TestUtils.scryRenderedDOMComponentsWithTag(about, 'p');
		expect(paragraphs.length).toEqual(2);
		expect(paragraphs[0].textContent).toEqual('foo');
		expect(paragraphs[1].textContent).toEqual('bar');
	});

	it('should switch language correctly', () => {
		about.switchLanguage('english');
		expect(about.state.content).toEqual(engCont);

		about.switchLanguage('foobar'); // shouldnt change
		expect(about.state.content).toEqual(engCont);
	});

	it('should fix header', () => {
		expect(about.state.hideFixHint).toEqual(false);
		
		about.fixHeader();
		expect(about.state.hideFixHint).toEqual(true);
		expect(fixHeader).toBeCalled();
	});
});