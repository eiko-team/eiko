if ("serviceWorker" in navigator) {
    if (!navigator.serviceWorker.controller) {
        navigator.serviceWorker.register("/eiko-sw.js");
    }
}

if (!isTokenValid(getCookie("Token"))) {
    window.location.replace("/login.html");
}

function logout() {
    log("logout");
    deleteCookie("token");
    window.location.replace("/login.html");
}

window.addEventListener("keydown", function(e) {
    if ((e.key === "Escape" || e.key === "Esc" || e.keyCode === 27) &&
        (e.target.nodeName === "BODY")) {
        closeNav();
        e.preventDefault();
        return false;
    }
}, true);