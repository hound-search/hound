const path = require('path');

module.exports = {
    entry: {
        'ui/js/hound.js': './ui/assets/js/hound.jsx',
        'ui/js/excluded_files.js': './ui/assets/js/excluded_files.jsx',
    },
    module: {
        rules: [
            {
                test: /\.jsx?$/,
                exclude: /node_modules/,
                use: {
                    loader: "babel-loader"
                }
            },
        ]
    },
    output: {
        filename: '[name]',
        path: path.resolve(__dirname, '.build')
    },
    resolve: {
        extensions: ['.js', '.jsx'],
    }
};
