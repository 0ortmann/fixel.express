const path = require('path');
const MiniCssExtractPlugin = require('mini-css-extract-plugin');
const webpack = require('webpack');

module.exports = {
	mode: 'production',
	devtool: 'source-map',
	entry: [
		'./src/entry.jsx'
	],
	output: {
		filename: 'bundle.js',
		path: path.resolve(__dirname, 'public/assets'),
		publicPath: '/assets/'
	},
	optimization: {
		minimize: true
	},
	plugins: [
		new MiniCssExtractPlugin({
			filename: 'bundle.css',
			chunkFilename: '[id].css'
		}),
		new webpack.DefinePlugin({
			'process.env': {
				'NODE_ENV': 'production'
			}
		})
	],
	module: {
		rules: [
		{
			test: /\.js$|\.jsx$/,
			exclude: /node_modules/,
			loader: 'babel-loader'
		},
		{
			test: /\.scss$/,
			use: [
				MiniCssExtractPlugin.loader,
				'css-loader',
				'sass-loader'
			],
			include: path.join(__dirname, 'src')
		},
		{
			test: /\.json$/,
			loader: 'json-loader',
			include: path.join(__dirname, 'src/content')
		}]
	}
};