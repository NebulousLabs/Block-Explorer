'use strict';

app.controller('MainCtrl', function($scope, ExplorerService, ConversionService){
  $scope.ConvSrvc = ConversionService;

  ExplorerService.getStatus()
      .success(function(data) {
          $scope.status = data;
      });
});
