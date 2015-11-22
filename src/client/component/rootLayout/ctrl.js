define('rootLayout/ctrl', [
    'cpStyle',
    'text!rootLayout/style.css',
    'text!rootLayout/tpl.html'
], function(
    cpStyle,
    style,
    tpl
){
    cpStyle(style);

    return function(ngModule){
        ngModule
            .directive('cpRootLayout', function(){
                return {
                    restrict: 'E',
                    template: tpl,
                    scope: {},
                    controller: ['$scope', function($scope){
                    }]
                };
            });
    }
});