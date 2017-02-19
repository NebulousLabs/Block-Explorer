'use strict';

app.directive('hostTable', function() {
    var linkFunc = function(scope, element, attrs) {
        scope.sortOrder = '';

        scope.sortBy = function(columnName) {
            if (scope.sortOrder == columnName) {
                scope.sortOrder = '-' + scope.sortOrder;
            } else {
                scope.sortOrder = columnName;
            }
        };
    };

    return {
        templateUrl: 'views/host_table.html',
        link: linkFunc
    }
});

