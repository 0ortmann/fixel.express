import React from 'react';
import { Route, IndexRoute } from 'react-router';

import App from './src/containers/App.jsx';
import About from './src/containers/About.jsx';

import ConnectK from './src/containers/ConnectK.jsx';


const routes = (
	<Route path='/' component={App}>
		<IndexRoute component={About} />
		<Route path='ai' component={ConnectK} />
	</Route>
);

export default routes;