define('3ditor/ctrl', [
    'cpStyle',
    'text!3ditor/style.css',
    'text!3ditor/tpl.html'
], function(
    cpStyle,
    style,
    tpl
){
    cpStyle(style);

    return function(ngModule){
        ngModule
            .directive('cp3ditor', function(){
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