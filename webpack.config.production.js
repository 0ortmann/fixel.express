var path = require('path');
var ExtractTextPlugin = require('extract-text-webpack-plugin');
var webpack = require('webpack');

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
		new ExtractTextPlugin('bundle.css'),
		new webpack.optimize.UglifyJsPlugin(),
		new webpack.optimize.DedupePlugin(),
		new webpack.DefinePlugin({
			'process.env': {
				'NODE_ENV': JSON.stringify(process.env.NODE_ENV)
			}
		})
	],
	module: {
		loaders: [{
			test: /\.js$|\.jsx$/,
			exclude: /node_modules/,
			loader: 'babel'
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