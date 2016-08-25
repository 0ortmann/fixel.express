jest.unmock('../src/components/main/Main.jsx');

import React from 'react';
import ReactDOM from 'react-dom';
import TestUtils from 'react-addons-test-utils';
import Main from '../src/components/main/Main.jsx';

describe('Main', () => {

	const content = {
		german: { description: ['foo', 'bar'], title: 't1'},
		english: { description: ['ping', 'pong'], title: 't2'}
	};

	let main;

	beforeEach(function() {
        main = TestUtils.renderIntoDocument( <Main content={content}/> );
    });

	it('should render', () => {
        const mainNode = ReactDOM.findDOMNode(main);
        expect(mainNode).toBeDefined();
        const paragraphs = TestUtils.scryRenderedDOMComponentsWithTag(main, 'p');
        expect(paragraphs.length).toEqual(2);
        expect(paragraphs[0].textContent).toEqual('foo')
        expect(paragraphs[1].textContent).toEqual('bar')
	});
});