'use strict';

app.controller('HostsCtrl', ['$scope', 'ExplorerService', function($scope, ExplorerService){
    ExplorerService.getHosts()
        .success(function(data){
            $scope.hosts = data.Hosts;
        });
}]);
