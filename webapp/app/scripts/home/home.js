'use strict';

angular
  .module('anvilmgr.home', [
    'ngResource',
    'ngRoute',
    'ui.router'
  ])
  .config(['$stateProvider', function($stateProvider) {

      var home = {
          name: 'app.home',
          url: '/',
          views: {
            '@': {
              templateUrl: 'views/home.html',
              controller: 'HomeController'
            }
          }
      };

      $stateProvider.state(home);

  }])
  .controller('HomeController', ['$scope', '$resource', function ($scope, $resource) {


    var Repos = $resource('/api/repo');

    $scope.repos = Repos.query();

  }]);

