import React from 'react';
import { Route, IndexRoute } from 'react-router';

import App from './src/containers/App.jsx';
import AboutMe from './src/components/about/AboutMe.jsx';

export default function configureRoutes(routeConfig) {
	const routes = routeConfig || [];
	return (
		<Route path='/' component={App}>
			<IndexRoute component={AboutMe} />
			{routes.map(r => {
				return <Route key={r.name} path={r.name} component={r.component} />;
			})}
		</Route>
	);
}