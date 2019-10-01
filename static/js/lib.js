/**
 * Make a POST request to the api
 * @param {string} url The api url
 * @param {object} body Data to send
 * @param {function} successCallback Function called on success
 * @param {function} failCallback Function called on failure
 */
function POST(url, body, successCallback = (e) => {},
    failCallback = (e) => {}) {
    return fetch("/api" + url, {
            method: "POST",
            body: JSON.stringify(body),
            headers: { "Content-Type": "application/json" }
        })
        .then((e) => { return e.json().then(successCallback); })
        .catch(failCallback);
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
    if (Notification.permission !== "granted") {
        Notification.requestPermission();
    } else {
        var notification = new Notification(title, { icon, body });
        notification.onclick = (e) => { self.open(onClickURL); };
    }
}

/**
 * update cookie value
 * @param {string} name name of the cookie to set
 * @param {string} value value to set the cookie to
 * @param {number} days number of days of validity of the cookie
 */
function createCookie(name, value, days = null) {
    var expire = "";
    if (days) {
        var d = new Date();
        d.setTime(d.getTime() + days * 24 * 60 * 60 * 1000);
        expire = `; expires=${d.toGMTString()}`;
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
    POST("/log", { user_token: getCookie("token"), message: msg });
}


function logout() {
    log("logout");
    deleteCookie("token");
    window.location.replace("/login.html");
}

/**
 * check is the token is valid
 */
function isTokenValid() {
    return getCookie("token") !== null;
}

/**
 * set style of all element(s) of a certain class
 * @param {string} className name of the class
 * @param {string} style style to apply on all elements
 */
function setStyleByClass(className, style) {
    var elt = document.getElementsByClassName(className);
    for (var i = 0; i < elt.length; i++) {
        elt[i].style = style;
    }
}

function setStyleByID(id, style) {
    document.getElementById(id).style = style;
}

function checkPassword(password) {
    return POST("/verify/password", { password }, (e) => {
        createCookie("pass_score", e.strength);
    });
}

function openNav() {
    log("openNav");
    setStyleByID("mySidenav", "display: block;");
}

function closeNav() {
    log("closeNav");
    setStyleByID("mySidenav", "display: none;");
}

function addlist(list) {
    let li = document.createElement("li");
    var uri = encodeURI(`/l/${list.id}`);
    li.innerHTML = `<a href="${uri}"><i class="material-icons">remove</i>${list.name}</a>`;
    var lists = document.getElementById("dropdown-lists");
    var last = lists.children[lists.children.length - 1]
    lists.appendChild(li);
    lists.appendChild(last);
}

function loadList() {
    log("loadList");
    POST("/list/getall", {}, (e) => {
        localStorage.setItem("lists", JSON.stringify(e));
        e.lists.forEach(addlist);
    }, (e) => {
        var lists = localStorage.getItem("lists");
        var json = { lists: [] };
        if (lists !== null) {
            json = JSON.parse(lists);
        }
        json.forEach(addlist);
    });
}

function createList(name = "Liste de course") {
    log("createList=" + name);
    var lists = localStorage.getItem("lists");
    var json = { lists: [] };
    if (lists !== null) {
        json = JSON.parse(lists);
    }
    json.lists.push(e);
    localStorage.setItem("lists", JSON.stringify(json));
    addlist(e.name);
    POST("/list/create", { name });
}

function shareList(id, email) {
    POST("/list/share", { id, email });
}

function deleteList(id) {
    POST("/list/delete", { id });
}
