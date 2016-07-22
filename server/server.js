import path from 'path'
import Express from 'express'
import webpack from 'webpack'
const config = require(path.resolve(__dirname, '../webpack.config.' + process.env.NODE_ENV + '.js'))

import React from 'react'
import { renderToString } from 'react-dom/server'

import App from '../src/components/App.jsx'

const compiler = webpack(config)
const app = new Express()
const port = 3000

if(process.env.NODE_ENV == 'development') {
    app.use(require('webpack-dev-middleware')(compiler, {
        publicPath: config.output.publicPath,
        noInfo: true 
    }));

    app.use(require('webpack-hot-middleware')(compiler));
    app.use(Express.static('public'))
}
else {
    app.use(Express.static('public'))
}


function handleRender(req, res) {
    var html = renderToString(<App />)

    res.send(renderFullPage(html))
}

function renderFullPage(html) {
    return `
        <!DOCTYPE HTML>
        <html>

        <head>
            <meta charset='utf-8'>
            <meta name="viewport" content="width=device-width, initial-scale = 1.0, user-scalable=yes" />
            <title>Fixel Express</title>
            <link rel='stylesheet' type='text/css' href='assets/bundle.css'>
            <link href='https://fonts.googleapis.com/css?family=Raleway' rel='stylesheet' type='text/css'>
            <link sizes="96x96" href="img/favicon.ico" type="image/png" rel="icon">
        </head>

        <body>
            <div id='content'>${html}</div>
            <script type='text/javascript' src='assets/bundle.js' charset='utf-8'></script>
        </body>

        </html>
    `;
}

app.get('/', handleRender);

var server = app.listen(3300, function() {
    var host = server.address().address
    var port = server.address().port

    console.log('App listening at http://%s:%s', host, port)
});
