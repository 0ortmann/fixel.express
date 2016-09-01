import 'babel-polyfill';

import React from 'react';
import { render } from 'react-dom';
import Root from './src/containers/Root.jsx';
import configureStore from './src/store/configureStore.js';

const store = configureStore()


render(
	<Root store={store} />,
	document.getElementById('content')
)
