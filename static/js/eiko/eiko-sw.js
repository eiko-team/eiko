// This is the service worker with the Cache-first network

var CACHE_NAME = "eiko-precache-v1.0.1";

function fromCache(request) {
    // Check to see if you have it in the cache
    // Return response
    // If not in the cache, then return
    return caches.open(CACHE_NAME).then(function(cache) {
        return cache.match(request).then(function(matching) {
            if (!matching || matching.status === 404) {
                return Promise.reject("no-match");
            }
            // see https://bugs.chromium.org/p/chromium/issues/detail?id=823392
            if (e.request.cache === 'only-if-cached' && e.request.mode !== 'same-origin') { return }
            return matching;
        });
    });
}

function updateCache(request, response) {
    return caches.open(CACHE_NAME).then(function(cache) {
        return cache.put(request, response);
    });
}

self.addEventListener("install", function(event) {
    self.skipWaiting();
    event.waitUntil(caches.open(CACHE_NAME).then(function(cache) {
        return cache.addAll([
            /* Add an array of files to precache for your app */
            "/",
            "/index.html",
            "/login.html",
            "/search/",
            "/favicon.ico",
            "/js/lib.js",
            "/js/search.js",
            "/js/eiko/eiko.js",
            "/js/login/login.js",
            "/js/color.js",
            "/css/eiko.css",
            "/img/loading.gif",
        ]);
    }));
});

// Allow sw to control of current page
self.addEventListener("activate", function(event) {
    // Remove previous cached data from disk.
    event.waitUntil(caches.keys().then(function(keyList) {
        return Promise.all(keyList.map(function(key) {
            if (key !== CACHE_NAME) { return caches.delete(key); }
        }));
    }));
    event.waitUntil(self.clients.claim());
});

// If any fetch fails, it will look for the request in the cache and serve it
// from there first
self.addEventListener("fetch", function(event) {
    if (event.request.method !== "GET") { return; }
    event.respondWith(fromCache(event.request).then(
        function(response) {
            // The response was found in the cache so we responde with it
            // and update the entry

            // This is where we call the server to get the newest version
            // of the file to use the next time we show view
            event.waitUntil(fetch(event.request).then(function(response) {
                return updateCache(event.request, response);
            }));
            return response;
        },
        function() {
            // The response was not found in the cache so we look for it on
            // the server
            return fetch(event.request).then(function(response) {
                // If request was success, add or update it in the cache
                event.waitUntil(updateCache(event.request, response.clone()));
                return response;
            }).catch(function(error) {});
        }
    ));
});

// for a 'add to home screen'
self.addEventListener("beforeinstallprompt", function(e) {
    // showInstallPromotion();
    log("beforeinstallprompt");
});

// log is the user has installed the app
self.addEventListener("appinstalled", function(e) {
    log("appinstalled");
});