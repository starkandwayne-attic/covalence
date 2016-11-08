var covalence = angular.module('covalence', [
  'ngRoute',
  'edgeBundling',
  'forceGraph'
]).controller('NavController', ['$scope', '$location', function ($scope, $location) {
    console.log("HERE")
    $scope.isCurrentPath = function (path) {
      console.log(path)
      return $location.path() == path;
    };
  }]);