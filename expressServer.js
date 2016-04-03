var express = require('express');
var path = require('path');
var webpack = require('webpack');
var config = require(path.resolve(__dirname, 'webpack.config.' + process.env.NODE_ENV + '.js'));
var compiler = webpack(config);


var app = express();

if(process.env.NODE_ENV == 'development') {
    app.use(require('webpack-dev-middleware')(compiler, {
        publicPath: config.output.publicPath,
        noInfo: true 
    }));

    app.use(require('webpack-hot-middleware')(compiler));
    app.use(express.static('public'));
}
else {
    app.use(express.static('public'));
}

app.get('/', function(req, res) {
    res.sendFile(path.join(__dirname, 'index.html'));
});


var server = app.listen(3300, function() {
    var host = server.address().address;
    var port = server.address().port;

    console.log('App listening at http://%s:%s', host, port);
});
