/**
 * Make a POST request to the api
 * @param {string} url The api url
 * @param {object} body Data to send
 * @param {function} successCallback Function called on success
 * @param {function} failCallback Function called on failure
 */
function POST(url, body, successCallback = (e) => {},
    failCallback = (e) => {}) {
    console.log("post:", "/api" + url, body);
    return fetch("/api" + url, {
            method: "POST",
            body: JSON.stringify(body),
            headers: { "Content-Type": "application/json" }
        })
        .then((e) => { return e.json().then(successCallback); })
        .catch(failCallback);
}

/**
 * Make a GET request to the server
 * @param {string} url The url
 * @param {object} body Data to send
 * @param {function} successCallback Function called on success
 * @param {function} failCallback Function called on failure
 */
function GET(url, successCallback = (e) => {},
    failCallback = (e) => {}) {
    return fetch(url, {
            method: "GET",
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
    li.id = list.id
    li.innerHTML = `<i class="material-icons">remove</i>${list.name}`;
    li.addEventListener("click", function(event) {
        createCookie("ListID", list.id);
        window.location.replace("/l/" + list.id);
    });
    var lists = document.getElementById("dropdown-lists");
    var last = lists.children[lists.children.length - 1];
    lists.appendChild(li);
    lists.appendChild(last);
}

/**
 * Removes all lists from sidenav
 */
function removeLists() {
    var lists = document.getElementById("dropdown-lists");
    var first = lists.children[0];
    var last = lists.children[lists.children.length - 1];
    while (lists.firstChild) {
        lists.removeChild(lists.firstChild);
    }
    lists.appendChild(first);
    lists.appendChild(last);
}

/**
 * load lists of the user. If the api fails, load lists from local storage
 */
function loadLists() {
    log("loadLists");
    var json = JSON.parse(localStorage.getItem("lists"));
    if (json === null || json.error === undefined) { json = []; }
    json.forEach(addlist);
    POST("/list/getall", {}, (e) => {
        if (e.error !== undefined) { e = [] }
        localStorage.setItem("lists", JSON.stringify(e));
        removeLists();
        e.forEach(addlist);
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
    json.push();
    localStorage.setItem("lists", JSON.stringify(json));
    addlist(name);
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
    var listID = getCookie("ListID");
    if (listID !== null && listID !== "" && listID !== "0") {
        return Number(listID);
    }
    var elts = document.URL.split("/");
    listID = elts[elts.length - 1];
    var fragmentIndex = listID.indexOf('#');
    if (fragmentIndex !== -1) {
        listID = listID.substring(0, fragmentIndex);
    }
    createCookie("ListID", listID);
    return Number(listID);
}

function hideConsumable(consumableId) {
    if (consumableId === undefined) { return; }
    var consumables = document.querySelector("tbody");
    var consumablesChilds = consumables.children;
    var consumable = null;
    var consumableIdStr = consumableId.toString()
    for (var i = 0; i < consumablesChilds.length; i++) {
        if (consumablesChilds[i].id.toString() === consumableIdStr) {
            consumable = consumablesChilds[i];
            break;
        }
    }
    if (consumable === null) { return; }
    consumables.removeChild(consumable);
}

function idToConsumable(consumableId) {
    var consumables = localStorage.getItem("consumables");
    if (consumables === null) { return; }
    return JSON.parse(consumables).filter(function(element) {
        return element.ID === consumableId;
    })[0];
}

function toggleDoneConsumable(consumableId) {
    var consumables = localStorage.getItem("consumables");
    var json = [];
    if (consumables !== null) {
        json = JSON.parse(consumables);
    }
    cons = json.filter(function(element) {
        return element.ID === consumableId;
    })
    if (cons.length === 0) { return; }
    cons = cons[0];
    cons.Done = !cons.Done;
    json = json.filter(function(element) {
        return element.ID !== consumableId;
    })
    json.push(cons);
    localStorage.setItem("consumables", JSON.stringify(json));
}

/**
 * create and return a new Unic IDentifier.
 */
function getNewUID() {
    var uid = Number(localStorage.getItem("UID")) + 1;
    localStorage.setItem("UID", uid);
    return uid;
}

function validateConsumable(consumableId) {
    return function(event) {
        toggleDoneConsumable(consumableId);
        fillConsumables();
        // TODO: send update to api
    }
}

/**
 * add a consumable to the page
 * @param {object} consumable to display
 */
function showConsumable(consumable) {
    console.log("showConsumable", consumable)
    if (!"content" in document.createElement("template")) { return; }
    var template = document.querySelector("#consumable");
    var clone = document.importNode(template.content, true);
    var td = clone.querySelectorAll("td");
    var tr = clone.querySelector("tr");
    tr.id = consumable.ID
    td[0].addEventListener("click", validateConsumable(consumable.ID));
    if (consumable.Done) {
        tr.classList.add("done");
        td[0].innerHTML = "<i class='material-icons'>radio_button_checked</i>";
    } else {
        td[0].innerHTML = "<i class='material-icons'>radio_button_unchecked</i>";
    }
    if (consumable.mode === "personnal") {
        td[1].textContent = consumable.name;
    } else {
        td[1].textContent = consumable.consumable.name;
        td[2].firstElementChild.classList.add("dot-green")
        td[3].firstElementChild.classList.add("dot-green")
        td[4].firstElementChild.classList.add("dot-green")
        // td[5].textContent = consumable.stock.pack_price / 100 + "€";
    }
    document.querySelector("tbody").appendChild(clone);
}

var loadingGifTimeoutID = [];

function showLoadingGif(status = false, timout = 2000, id = "main") {
    if (status) {
        loadingGifTimeoutID.push(setTimeout(function(e) {
            document.getElementById(id + "-loading-gif").style.display = "";
        }, timout))
    } else {
        loadingGifTimeoutID.forEach(clearTimeout);
        loadingGifTimeoutID = [];
        document.getElementById(id + "-loading-gif").style.display = "none";
    }
}

function addConsumableLocalStorage(consumable) {
    var consumables = localStorage.getItem("consumables");
    var json = [];
    if (consumables !== null) {
        json = JSON.parse(consumables);
    }

    json = json.filter(function(element) {
        return element.ID !== consumable.ID;
    });
    json.push(consumable)
    localStorage.setItem("consumables", JSON.stringify(json));
}

function addConsumablesFromList(consumables) {
    if (consumables.error !== undefined) { return; }
    consumables.forEach(addConsumableLocalStorage);
    fillConsumables();
}

function fetchConsumables() {
    var current = getCurrentListID();
    POST("/list/getcontent", { ID: current }, addConsumablesFromList);
    var lists = JSON.parse(localStorage.getItem("lists"));
    if (lists === null) { return; }
    for (var i = 0; i < lists.length; i++) {
        var ID = lists[i].id;
        if (ID === current) { continue; }
        POST("/list/getcontent", { ID }, addConsumablesFromList);
    }
}

/**
 * add all consumables to the page according to the list id from local storage.
 * @param {boolean} fetch updated version on the server
 */
function fillConsumables(fetch = false) {
    var list = { id: getCurrentListID() };
    var consumables = localStorage.getItem("consumables");
    var json = [];
    if (consumables !== null) {
        json = JSON.parse(consumables);
    }
    if (list.id === 0) { return; };
    showLoadingGif(true);
    var done = [];
    json.sort(function(a, b) {
        return a.ID - b.ID;
    });
    json.forEach(function(element) {
        console.log(element, element.list_id, list.id)
        if (element.list_id === list.id) {
            hideConsumable(element.ID);
            if (element.Done) {
                done.push(element);
            } else {
                showConsumable(element);
            }
        }
    });
    done.forEach(showConsumable);
    showLoadingGif(false);
    if (fetch) {
        fetchConsumables();
    }
}

function getTargetPosition() {
    if (getCookie("posMode") === "local") {
        if (!navigator.geolocation) {
            // cannot use localisation *sad*
            return;
        }
        navigator.geolocation.getCurrentPosition(function(e) {
            createCookie("posLat", e.coords.latitude);
            createCookie("posLon", e.coords.longitude);
            createCookie("posAcc", e.coords.accuracy);
        }, function(e) {}, options = {
            enableHighAccuracy: true,
            timeout: 5000,
            maximumAge: 0
        });
    }
}

function goBackList(fragment = "") {
    // removing Events
    if (fragment.type !== undefined) { fragment = ""; }
    var listID = getCookie("ListID");
    if (listID === "" || listID === "0") {
        window.location.replace("/" + fragment);
    } else {
        window.location.replace("/l/" + listID + fragment);
    }
}