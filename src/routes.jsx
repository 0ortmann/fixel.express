import React from 'react';
import { IndexRoute, Route } from 'react-router';

import App from './containers/App.jsx';
import AboutMe from './components/about/AboutMe.jsx';
import Imprint from './components/imprint/Imprint.jsx';

export default function configureRoutes() {
	return (
		<Route path='/' component={App}>
			<IndexRoute component={AboutMe} />
            <Route key='imprint' path='imprint' component={Imprint} />
		</Route>
	);
}