var covalence = angular.module('covalence', []);

covalence.controller('VMsFilterController', function VMsFilterController($scope, $http) {
    $http({
        method: 'GET',
        url: '/vms'
    }).then(function successCallback(response) {
        var deployments = []
        $scope.vms = response.data
        for (var i = 0; i < $scope.vms.length; i++) {
            deployments.push($scope.vms[i].deployment_name)
        }
        $scope.deployments = deployments.filter((v, i, a) => a.indexOf(v) === i);
    }, function errorCallback(response) {
        // called asynchronously if an error occurs
        // or server returns response with an error status.
    });
});