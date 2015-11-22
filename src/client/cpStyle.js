define('cpStyle', [
    'text!topLevelCss'
], function(
    topLevelCss
){

    var id = 'cp-style-id';

    var el = document.createElement('style');
    el.id = id;
    document.head.appendChild(el);

    function cpStyle(styleTxt){

        if(typeof styleTxt !== 'string') {
            throw 'styleTxt must be a string';
        }

        el.innerHTML += styleTxt + ' ';

    };

    cpStyle(topLevelCss);

    return cpStyle;
});
