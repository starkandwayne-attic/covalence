angular.
  module('covalence').
  config(['$locationProvider', '$routeProvider',
    function config($locationProvider, $routeProvider) {
      $locationProvider.hashPrefix('!');

      $routeProvider.
        when('/edge', {
          template: '<edge-bundling></edge-bundling>'
        }).
        when('/force', {
          template: '<force-graph></force-graph>'
        }).
        otherwise('/edge');
    }
  ]);