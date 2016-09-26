import React from 'react';
import { Route, IndexRoute } from 'react-router';
import path from 'path';

import App from './src/containers/App.jsx';
import AboutMe from './src/components/about/AboutMe.jsx';

import Game from './src/components/connect-k/Game.jsx';

const config = require('./routes.config.' + process.env.NODE_ENV + '.js');

const routes = (
	<Route path='/' component={App}>
		<IndexRoute component={AboutMe} />
		{config.routes.map(r => {
			return <Route key={r.name} path={r.name} component={r.component} />
		})}
	</Route>
);

export default routes;