var covalence = angular.module('covalence', [
  'ngRoute',
  'edgeBundling',
  'forceGraph'
]).controller('NavController', ['$scope', '$location', function ($scope, $location) {
    $scope.isCurrentPath = function (path) {
      return $location.path() == path;
    };
  }]);