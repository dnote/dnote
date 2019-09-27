/* eslint-disable no-restricted-globals, func-names, no-var, prefer-template */

self.addEventListener('install', function(event) {
  var offlineRequest = new Request('static/offline.html');

  event.waitUntil(
    fetch(offlineRequest).then(function(response) {
      return caches.open('offline').then(function(cache) {
        console.log('cached offline page', response.url);
        return cache.put(offlineRequest, response);
      });
    })
  );
});

self.addEventListener('fetch', function(event) {
  var request = event.request;
  if (request.method === 'GET') {
    event.respondWith(
      fetch(request).catch(function(error) {
        console.error(
          'onfetch Failed. Serving cached offline fallback ' + error
        );
        return caches.open('offline').then(function(cache) {
          return cache.match('offline.html');
        });
      })
    );
  }
});
