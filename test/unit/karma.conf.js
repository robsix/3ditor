module.exports = function(config) {
    config.set({

        basePath: '../../',

        frameworks: ['jasmine', 'requirejs'],

        files: [
            {pattern: 'test/unit/require.conf.js', included: true},
            {pattern: 'test/unit/tests/global.js', included: true},
            {pattern: 'src/client/**/*.*', included: false},
            {pattern: 'test/unit/tests/**/*.*', included: false},
        ],

        plugins: [
            'karma-jasmine',
            'karma-requirejs',
            'karma-coverage',
            'karma-html-reporter',
            'karma-phantomjs-launcher',
            'karma-chrome-launcher',
            'karma-firefox-launcher',
            'karma-safari-launcher',
            'karma-ie-launcher'
        ],

        reporters: ['coverage', 'html', 'progress'],

        preprocessors: {
            'src/client/component/**/*.js': ['coverage'],
            'src/client/service/**/*.js': ['coverage']
        },

        coverageReporter: {
            type: 'html',
            dir: 'test/unit/coverage/',
            includeAllSources: true
        },

        htmlReporter: {
            outputDir: 'results' //it is annoying that this file path isn't from basePath :(
        },

        colors: true,

        logLevel: config.LOG_INFO,

        autoWatch: false,

        browsers: ['Chrome'/*, 'PhantomJS', 'Firefox', 'IE', 'Safari'*/],

        captureTimeout: 5000,

        singleRun: true
    });
};