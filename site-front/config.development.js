import Game from './src/components/connect-k/Game.jsx';

module.exports = {
	routes: [ { name: 'ai', component: Game } ],
	api: {
		resourcesHost: 'http://localhost:3300',
		connectKHost: 'http://localhost:5000'
	},
	logging: true
};