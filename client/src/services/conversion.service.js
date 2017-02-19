'use strict';

app.factory('ConversionService', function(){
    var conversionSrvc = {};

    conversionSrvc.hastingsToString = function(hastings) {
        return (hastings / Math.pow(10, 24)).toFixed(2);
    }

    conversionSrvc.bytesToGB = function(bytes) {
        return (bytes / Math.pow(10, 9)).toFixed(2);
    }

    conversionSrvc.bytesToTB = function(bytes) {
        return (bytes / Math.pow(10, 12)).toFixed(2);
    }

    return conversionSrvc;
});
