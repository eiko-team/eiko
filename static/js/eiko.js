// console.log = function() {}

if ("serviceWorker" in navigator) {
    if (navigator.serviceWorker.controller) {
        console.log("[SW] Active service worker found, no need to register");
    } else {
        // Register the service worker
        navigator.serviceWorker.register("/eiko-sw.js", { scope: "./" })
            .then(function(reg) {
                console.log("[SW] Service worker registered for scope: " +
                    reg.scope);
            });
    }
}