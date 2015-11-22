define('registry', [
    'ng',
    //services

    //components
    '3ditor/ctrl',
    'rootLayout/ctrl'
], function(
    ng,
    //services

    //components
    editorCtrl,
    rootLayoutCtrl
){
    var registry = ng.module('cp.registry', []);
    [].slice.call(arguments).forEach(function(arg){
        if(arg !== ng){
            arg(registry);
        }
    });

    return registry;
});