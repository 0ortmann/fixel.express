var webpack = require('webpack');
var path = require('path');

module.exports = {
	devtool: 'cheap-module-eval-source-map',
	entry: [
		'webpack-hot-middleware/client',
		'./entry.jsx'
	],
	plugins: [
		new webpack.optimize.OccurrenceOrderPlugin(),
		new webpack.HotModuleReplacementPlugin(),
		new webpack.DefinePlugin({
			'process.env': {
				'NODE_ENV': JSON.stringify(process.env.NODE_ENV)
			}
		})
	],
	output: {
		filename: 'bundle.js',
		path: path.resolve(__dirname, 'public/assets'),
		publicPath: '/assets/'
	},
	module: {
		loaders: [{
			test: /\.js$|\.jsx$/,
			exclude: /node_modules/,
			loader: 'babel'
		},
		{
			test: /\.scss$/,
			loader: 'style-loader!css-loader!sass-loader',
			include: path.join(__dirname, 'src')
		},
		{
			test: /\.json$/,
			loader: 'json-loader',
			include: path.join(__dirname, 'src/content')
		}]
	}
};