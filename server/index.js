require('babel-core/register');
require.extensions['.scss'] = () => {
	return;
};
require('./server.js');