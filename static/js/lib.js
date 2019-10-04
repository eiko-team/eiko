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

/**
 * set style of an element of a certain id
 * @param {string} id name of the id
 * @param {string} style style to apply on the element
 */
function setStyleByID(id, style) {
    document.getElementById(id).style = style;
}

/**
 * set the cookie "pass_score" with the score of the password to enable password
 * strenght display
 * @param {string} password to check
 */
function checkPassword(password) {
    return POST("/verify/password", { password }, (e) => {
        createCookie("pass_score", e.strength);
    });
}

/**
 * open thena bar
 */
function openNav() {
    log("openNav");
    setStyleByID("mySidenav", "display: block;");
}

/**
 * close thena bar
 */
function closeNav() {
    log("closeNav");
    setStyleByID("mySidenav", "display: none;");
}

/**
 * add a list in the nav bar
 * @param {object} list to add
 */
function addlist(list) {
    let li = document.createElement("li");
    var uri = encodeURI(`/l/${list.id}`);
    li.innerHTML = `<a href="${uri}"><i class="material-icons">remove</i>${list.name}</a>`;
    var lists = document.getElementById("dropdown-lists");
    var last = lists.children[lists.children.length - 1];
    lists.appendChild(li);
    lists.appendChild(last);
}

/**
 * load lists of the user. If the api fails, load lists from local storage
 */
function loadLists() {
    log("loadLists");
    POST("/list/getall", {}, (e) => {
        localStorage.setItem("lists", JSON.stringify(e));
        e.lists.forEach(addlist);
    }, (e) => {
        var lists = localStorage.getItem("lists");
        var json = [];
        if (lists !== null) {
            json = JSON.parse(lists);
        }
        json.forEach(addlist);
    });
}

/**
 * create a list, sends it to the api and store it in local storage
 * @param {string} name of the list
 */
function createList(name = "Liste de course") {
    log("createList=" + name);
    var lists = localStorage.getItem("lists");
    var json = [];
    if (lists !== null) {
        json = JSON.parse(lists);
    }
    json.lists.push(e);
    localStorage.setItem("lists", JSON.stringify(json));
    addlist(e.name);
    POST("/list/create", { name });
}

/**
 * Shares a list with a user
 * @param {interger} id of the list
 * @param {string} email of the user to share the list with
 */
function shareList(id, email) {
    POST("/list/share", { id, email });
}

/**
 * Delete a list
 * @param {interger} id of the list
 */
function deleteList(id) {
    POST("/list/delete", { id });
}

/**
 * retrieve consumables of a list from the api and stores it in the local
 * storage
 * @param {object} list
 */
function getConsumables(list) {
    POST("/list/get", { list }, (e) => {
        var consumables = localStorage.getItem("consumables");
        var json = [];
        if (consumables !== null) {
            json = JSON.parse(consumables);
        }
        json = json.filter(function(element) {
            return element.list_id !== list.ID;
        })
        e.forEach(json.push);
        localStorage.setItem("consumables", JSON.stringify(json));
    });
}

/**
 * return the id of the list displayed (or to be displayed) on the page and
 * stores it in the cookies
 */
function getCurrentListID() {
    var elts = document.URL.split("/");
    var ListID = Number(elts[elts.length - 1]);
    createCookie("ListID", ListID);
    return ListID;
}

/**
 * add a consumable to the page
 * @param {object} consumable to display
 */
function showConsumable(consumable) {
    if (!"content" in document.createElement("template")) { return; }
    var template = document.querySelector("#consumable");
    var clone = document.importNode(template.content, true);
    var td = clone.querySelectorAll("td");
    td[1].textContent = consumable.consumable.name;
    td[5].textContent = consumable.stock.pack_price / 100 + "€";
    document.querySelector("tbody").appendChild(clone);
}


/**
 * add all consumables to the page according to the list id.
 * if the api fails, uses local storage.
 * else store the result in local storage
 * @param {object} consumable to display
 */
function fillConsumables() {
    log("fillConsumables");
    var list = {id: getCurrentListID()};
    POST("/list/get", { list }, function(e) {
        e.forEach(showConsumable);
        getConsumables(list);
    }, function(e) {
        var consumables = localStorage.getItem("consumables");
        var json = [];
        if (consumables !== null) {
            json = JSON.parse(consumables);
        }
        json.forEach(function(element) {
            if (element.list_id === list.ID) {
                showConsumable(element);
            }
        });
    })
}