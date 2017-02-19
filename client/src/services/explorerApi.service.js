'use strict';

app.factory('ExplorerService', ['$http', function($http){
    var baseurl = '/api';

    var explorerService = {};

    explorerService.getHosts = function(){
        return $http.get(baseurl + '/hosts/')
    };

    explorerService.getStatus = function(){
        return $http.get(baseurl + '/status/')
    };

    explorerService.getBlockHash = function(hash){
        return $http.get(baseurl + '/block/hash/' + hash)
    };

    return explorerService;
}]);
