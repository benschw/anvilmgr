'use strict';
console.log('loaded');
/**
 * @ngdoc overview
 * @name anvilmgrApp
 * @description
 * # satisApp
 *
 * Main module of the application.
 */
angular
  .module('anvilmgrApp', [
    'ngAnimate',
    'ngCookies',
    'ngResource',
    'ngRoute',
    'ngSanitize',
    'ngTouch',
    'ui.router',
    'ui.bootstrap',
    'anvilmgr.home',
    'anvilmgr.api'
  ])
  .run(['$rootScope', '$state', '$stateParams',
    function ($rootScope, $state, $stateParams) {
      console.log('app');
      $rootScope.$state = $state;
      $rootScope.$stateParams = $stateParams;
  }])
  .config(['$stateProvider', '$urlRouterProvider', function ($stateProvider, $urlRouterProvider) {
    // If the url is ever invalid, e.g. '/asdf', then redirect to '/' aka the home state
    $urlRouterProvider.otherwise('/');

    var app = {
      name: 'app',
      abstract: true,
      url: '',
      views: {
        'header': {
          templateUrl: 'views/header.html',
        },
        // '': {}, // skip in root state
        'footer': {
          templateUrl: 'views/footer.html'
        }          }
    };

    $stateProvider.state(app);

  }]);


