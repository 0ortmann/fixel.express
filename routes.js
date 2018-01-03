import React from 'react';
import { Route, IndexRoute } from 'react-router';

import App from './src/containers/App.jsx';
import AboutMe from './src/components/about/AboutMe.jsx';

export default function configureRoutes() {
	return (
		<Route path='/' component={App}>
			<IndexRoute component={AboutMe} />
		</Route>
	);
}