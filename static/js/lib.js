/**
 * Make a POST request to the api
 * @param {string} url The api url
 * @param {object} body Data to send
 * @param {function} successCallback Function called on success
 * @param {function} failCallback Function called on failure
 */
function POST(url, body, successCallback = (e) => {},
    failCallback = (e) => {}) {
    // if (typeof(failCallback) === typeof(async)) {
    //     async = failCallback
    //     failCallback = (e) => {}
    // }
    // var req = new XMLHttpRequest();
    // req.open("POST", "/api/"+url, async);
    // req.addEventListener("load", wrap(successCallback, req));
    // req.addEventListener("error", wrap(failCallback, req));
    // req.send(JSON.stringify(body));
    fetch("/api" + url, {
            method: "POST",
            body: JSON.stringify(body),
            headers: { 'Content-Type': 'application/json' }
        })
        .then((e) => { e.json().then(successCallback) })
        .catch(failCallback)
}

/**
 * Send a notification to the user if the user has granted permissions
 * It might be usefull to use Server-sent events
 * https://en.wikipedia.org/wiki/Server-sent_events
 * @param {string} title title of the notification
 * @param {string} body message inside the notification
 * @param {string} onClickURL the URL to redirect the user if there is a click
 * @param {string} icon icon to display in the notification
 */
function notify(title, body, onClickURL = "#", icon = "/favicon.ico") {
    if (Notification.permission !== 'granted')
        Notification.requestPermission();
    else {
        var notification = new Notification(title, {
            icon: icon,
            body: body,
        });
        notification.onclick = function() {
            self.open(onClickURL);
        };
    }
}

/**
 * update cookie value
 * @param {string} name name of the cookie to set
 * @param {string} value value to set the cookie to
 * @param {number} days number of days of validity of the cookie
 */
function createCookie(name, value, days) {
    var expire = ""
    if (days) {
        var d = new Date();
        d.setTime(d.getTime() + days * 24 * 60 * 60 * 1000)
        expire = `; expires=${d.toGMTString()}`
    }
    document.cookie = `${name}=${value}${expire}; path=/`;
}

/**
 * make a cookie invalid, thus deleting it
 * @param {string} name name of the cookie to delete
 */
function deleteCookie(name) {
    createCookie(name, "", -1);
}

/**
 * get cookie value
 * @param {string} name name of the cookie
 */
function getCookie(name) {
    var re = new RegExp(name + "=([^;]+)");
    var value = re.exec(document.cookie);
    return (value != null) ? unescape(value[1]) : null;
}

/**
 * Send logs to backend
 * @param {object} msg Data to send
 */
function log(msg) {
    POST("/log", { user_token: getCookie("token"), message: msg })
}


/**
 * check is the token is valid
 */
function isTokenValid() {
    return getCookie("token") !== null
}

/**
 * set style of all element(s) of a certain class
 * @param {string} className name of the class
 * @param {string} style style to apply on all elements
 */
function setStyleByClass(className, style) {
    var elt = document.getElementsByClassName(className);
    for (var i = 0; i < elt.length; i++)
        elt[i].style = style
}

function setStyleByID(id, style) {
    document.getElementById(id).style = style
}

function checkPassword(password) {
    return {score: Math.floor((Math.random() * 5) + 1) - 1};
    return POST("/verify", {password: password}, (e) => {
        return e
    })
}