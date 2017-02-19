'use strict';

app.directive('hostRow', function() {
    var linkFunction = function (scope, element) {
        scope.getTotalStorage = function(){
            if (scope.host.TotalStorage > Math.pow(10, 12)) {
                return (scope.host.TotalStorage / Math.pow(10, 12)) + ' TB';
            } else {
                return (scope.host.TotalStorage / Math.pow(10, 9)) + ' GB';
            }
        }

        scope.getSiaPerGBPerMonth = function(){
            return (scope.host.Price) / (4.32 * Math.pow(10, 12));
        }
    }
    return {
        templateUrl: 'views/host_row.html',
        link: linkFunction,
    }
});

