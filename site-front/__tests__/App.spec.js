import React from 'react';
import ReactDOM from 'react-dom';
import TestUtils from 'react-addons-test-utils';

import { App } from '../src/containers/App.jsx';

describe('App', () => {

	const langProps = { 
		description: ['foo', 'bar'],
		title: 't1',
		box: 'box'
	}; 

	const getLanguage = jest.genMockFunction();


	it('should render', () => {
		const app = TestUtils.renderIntoDocument(
			<App langProps={langProps} getLanguage={getLanguage}>
				<div className='CHILD' />
			</App>
		);
		const appNode = ReactDOM.findDOMNode(app);
		expect(appNode).toBeDefined();

		expect(getLanguage).toBeCalledWith('german');
		
		const child = TestUtils.findRenderedDOMComponentWithClass(app, 'CHILD');
		expect(child).toBeDefined();
	});
});