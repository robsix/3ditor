var page = require('./page.js');

'use strict';

describe('global', function() {

    beforeEach(function(){
        page.get();
    });

    it('page title should be "3ditor"', function() {
        expect(browser.getTitle()).toEqual('3ditor');
    });

});
