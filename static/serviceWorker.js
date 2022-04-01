const addResourcesToCache = async (resources) => {
    const cache = await caches.open('v1');
    await cache.addAll(resources);
};

self.addEventListener('install', (event) => {
    event.waitUntil(
        addResourcesToCache([
            '/'
        ])
    );
});

self.addEventListener('fetch', function (event) {
    event.respondWith(
        fetch(event.request)
            .then(function (response) {
                caches.open('v1').then(function (cache) {
                    cache.put(event.request, response.clone())
                })
                return response.clone()
            })
            .catch(function () {
                return caches.match(event.request);
            })
    );
});
