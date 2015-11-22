var tests = [];
for (var file in window.__karma__.files) {
    if (window.__karma__.files.hasOwnProperty(file)) {
        if (/test\/unit\/tests\/component\/.*\.js$/.test(file) || /test\/unit\/tests\/service\/.*\.js$/.test(file)) {
            tests.push(file);
        }
    }
}

requirejs.config({
    baseUrl: '/base/src/client/component',
    paths: {
        ng: '../lib/angular/angular',
        ngRoute: '../lib/angular-route/angular-route',
        three: '../lib/three/index',
        moment: '../lib/moment/min/moment-with-locales',
        text: '../lib/requirejs-text/text',
        registry: '../registry',
        service: '../service',
        topLevelCss: '../topLevel.css',
        ngMock: '../lib/angular-mocks/angular-mocks'
    },
    shim: {
        ng: {
            exports: 'angular'
        },
        ngMock:{
            deps: ['ng']
        },
        three: {
            exports: 'THREE'
        }
    },
    deps: tests,
    callback: window.__karma__.start
});