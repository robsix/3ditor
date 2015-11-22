require([
    'ng',
    'ngRoute',
    'registry'
], function(
    ng
){
    var app = ng.module('app', [
            'ngRoute',
            'cp.registry'
        ]);

    app.config(['$routeProvider', function ($routeProvider) {
        $routeProvider
            .when('/', {
                template: '<cp-3ditor class="cp-def cp-fill cp-dark-theme" ng-cloak></cp-3ditor>'
            })
            .otherwise({
                redirectTo: '/'
            });
    }]);

    ng.element(document).ready(function () {
        ng.bootstrap(document, ['app']);
    });
});