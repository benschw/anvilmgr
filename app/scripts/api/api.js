'use strict';

angular
  .module('anvilmgr.api', [
    'ngResource',
    'ngRoute',
    'ui.router'
  ])
  .config(['$stateProvider', function($stateProvider) {

      var home = {
          name: 'app.api',
          url: '/docs',
          views: {
            '@': {
              templateUrl: 'views/api.html',
              controller: 'ApiController'
            }
          }
      };

      $stateProvider.state(home);

  }])
  .controller('ApiController', ['$scope', '$resource', function ($scope, $resource) {





  }]);

