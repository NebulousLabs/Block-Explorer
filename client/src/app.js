'use strict';

var app = angular.module('siaExplorerApp', ['ngRoute']);

app.config(['$routeProvider', function($routeProvider) {
        $routeProvider
        .when('/', {
            controller: "MainCtrl",
            templateUrl: 'views/main.html',
        })
        .when('/hosts/', {
            controller: "HostsCtrl",
            templateUrl: 'views/hosts.html',
        })
        .when('/block/:blockHash', {
            controller: "BlockCtrl",
            templateUrl: 'views/block.html',
        })
        .otherwise({redirectTo: '/'});
    }]);

// Initialize foundation once the page has loaded.
app.run(function($timeout){
    $timeout(function() {
        $(document).foundation();
    }, 500);
});
