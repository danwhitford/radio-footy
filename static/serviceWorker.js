const addResourcesToCache = async (resources) => {
    const cache = await caches.open('v1');
    await cache.addAll(resources);
};

self.addEventListener('install', (event) => {
    event.waitUntil(
        addResourcesToCache([
            '/',
            '/index.html',
        ])
    );
});

self.addEventListener('fetch', function(event) {
    console.log(event.request)
    console.log(caches.match(event.request))
    event.respondWith(
      fetch(event.request).catch(function() {
        return caches.match(event.request);
      })
    );
  });
