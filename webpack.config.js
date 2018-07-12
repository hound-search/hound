const path = require('path');

module.exports = {
    entry: {
        'ui/js/hound.js': './ui/assets/js/hound.js',
        'ui/js/excluded_files.js': './ui/assets/js/excluded_files.js',
    },
    module: {
        rules: [
            {
                test: /\.js$/,
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
    }
};