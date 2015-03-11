'use strict';

angular
  .module('anvilmgr.api', [
    // 'ngResource',
    'ngRoute',
    'ui.router'
  ])
  .config(['$stateProvider', function($stateProvider) {

      var api = {
          name: 'app.api',
          url: '/docs',
          views: {
            '@': {
              templateUrl: 'views/api.html',
              controller: function(){
                console.log('docs controller')
              }
              // controller: 'ApiController'
            }
          }
      };

      $stateProvider.state(api);

  }]);
  // .controller('ApiController', ['$scope', '$resource', function ($scope, $resource) {





  // }]);

