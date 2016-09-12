import React from 'react';
import { Route, IndexRoute } from 'react-router';

import App from './src/containers/App.jsx';
import AboutMe from './src/components/about/AboutMe.jsx';

import Game from './src/components/connect-k/Game.jsx';


const routes = (
	<Route path='/' component={App}>
		<IndexRoute component={AboutMe} />
		<Route path='ai' component={Game} />
	</Route>
);

export default routes;