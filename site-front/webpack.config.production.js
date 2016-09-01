var path = require('path');
var ExtractTextPlugin = require('extract-text-webpack-plugin');

module.exports = {
	devtool: 'source-map',
	entry: [
		'./entry.jsx'
	],
	output: {
		filename: 'bundle.js',
		path: path.resolve(__dirname, 'public/assets'),
		publicPath: '/assets/'
	},
	plugins: [
		new ExtractTextPlugin('bundle.css')
	],
	module: {
		loaders: [{
			test: /\.js$|\.jsx$/,
			exclude: /node_modules/,
			loaders: ['babel-loader?presets[]=es2015&presets[]=react'],
		},
		{
			test: /\.scss$/,
			loader: ExtractTextPlugin.extract('style', 'css!sass'),
			include: path.join(__dirname, 'src')
		},
		{
			test: /\.json$/,
			loader: 'json-loader',
			include: path.join(__dirname, 'src/content')
		}]
	}
};