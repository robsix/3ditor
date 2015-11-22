require.config({
    baseUrl: 'component',
    paths: {
        ng: '../lib/angular/angular',
        ngRoute: '../lib/angular-route/angular-route',
        three: '../lib/three/index',
        moment: '../lib/moment/min/moment-with-locales',
        text: '../lib/requirejs-text/text',
        registry: '../registry',
        service: '../service',
        cpStyle: '../cpStyle',
        topLevelCss: '../topLevel.css'
    },
    shim: {
        ng: {
            exports: 'angular'
        },
        ngRoute: {
            deps: ['ng'],
            exports: 'angular'
        },
        three: {
            exports: 'THREE'
        },
        'moment': {
            exports: 'moment'
        }
    }
});