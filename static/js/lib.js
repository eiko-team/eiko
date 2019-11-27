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
 * get cookie value
 * @param {string} name name of the cookie
 */
function getCookie(name) {
    var re = new RegExp(name + "=([^;]+)");
    var value = re.exec(document.cookie);
    return (value != null) ? unescape(value[1]) : null;
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
    if (Notification.permission !== "granted" &&
        getCookie("notification") !== "no") {
        var promise = Notification.requestPermission();
        promise.then(function(event) {
            // onFullfiled
            notify(title, body, onClickURL, icon);
            createCookie("notification", "yes");
        }, function(event) {
            // onRejected
            M.toast({ html: "We won\'t send any notifications to you" })
            createCookie("notification", "no");
        })
    } else {
        var notification = new Notification(title, { icon, body });
        notification.onclick = (e) => { self.open(onClickURL); };
    }
}

/**
 * Send a toast to the user
 * @param {string} html content of the toast
 */
function toast(toast) {
    toast.classes = "rounded";
    M.toast(toast)
}

/**
 * make a cookie invalid, thus deleting it
 * @param {string} name name of the cookie to delete
 */
function deleteCookie(name) {
    createCookie(name, "", -1);
}

/**
 * Send logs to backend
 * @param {object} msg Data to send
 */
function log(msg) {
    POST("/log", { user_token: getCookie("token"), message: msg });
}

/**
 * logout the user
 */
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
 * returns the elt associated to the eltId in storageArea
 * @param {integer} eltId id of the consumable to return
 * @param {string} storageArea place inside the local storage
 * @returns {object} consumable associated to the eltId
 */
function idToElt(eltId, storageArea) {
    var consumables = JSON.parse(localStorage.getItem(storageArea));
    if (consumables === null) { return; }
    return consumables.filter(function(element) {
        return element.ID === eltId;
    })[0];
}

/**
 * replace e1 by e2 in local storage
 * @param {object} e1 elt to replace
 * @param {integer} e1.ID id to replace, it is the only value of c1 used
 * @param {object} e2 elt to take place
 * @param {string} storageArea place inside the local storage
 */
function updateLocalStorage(e1, e2, storageArea) {
    var json = JSON.parse(localStorage.getItem(storageArea));
    if (json === null) {
        json = [e2];
        localStorage.setItem(storageArea, JSON.stringify(json));
    }
    // removing e1
    json = json.filter(function(element) {
        return element.ID !== e1.ID;
    })
    // pushing e2
    json.push(e2);
    localStorage.setItem(storageArea, JSON.stringify(json));
}

/**
 * insert elt in local storage
 * @param {object} elt elt to insert
 * @param {string} storageArea place inside the local storage
 */
function insertLocalStorage(elt, storageArea) {
    var json = JSON.parse(localStorage.getItem(storageArea));
    if (json === null) { json = []; }
    json = json.filter(function(e) { return e.ID !== elt.ID; });
    json.push(elt);
    localStorage.setItem(storageArea, JSON.stringify(json));
}

/**
 * remove elt in local storage
 * @param {object} elt elt to remove
 * @param {string} storageArea place inside the local storage
 */
function removeLocalStorage(elt, storageArea) {
    var notifications = JSON.parse(localStorage.getItem(storageArea));
    if (notifications === null) { return; }
    notifications = notifications.filter(function(n) {
        return n.ID !== elt.ID
    })
    localStorage.setItem(storageArea, JSON.stringify(notifications));
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
 * add a list in the nav bar
 * @param {object} list to add
 */
function addlist(list) {
    let li = document.createElement("li");
    var uri = encodeURI(`/l/${list.id}`);
    li.id = list.id
    li.innerHTML = `<a><i class="material-icons">list_alt</i>${list.name}</a>`;
    li.addEventListener("click", function(event) {
        createCookie("ListID", list.id);
        window.location.replace("/l/" + list.id);
    });
    var lists = document.getElementById("collapsible-lists");
    var last = lists.children[lists.children.length - 1];
    lists.appendChild(li);
    lists.appendChild(last);
}

/**
 * Removes all lists from sidenav
 */
function removeLists() {
    var lists = document.getElementById("collapsible-lists");
    var last = lists.children[lists.children.length - 1];
    while (lists.firstChild) {
        lists.removeChild(lists.firstChild);
    }
    lists.appendChild(last);
}

/**
 * load lists from local storage and update local storage from api.
 */
function loadLists() {
    log("loadLists");
    var json = JSON.parse(localStorage.getItem("lists"));
    if (json === null || json.error !== undefined) { json = []; }
    json.forEach(addlist);
    POST("/list/getall", {}, (e) => {
        if (e === null || e.error !== undefined) { e = [] }
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
    var list = { name, id: getNewUID() };
    addlist(list);
    insertLocalStorage(list, "lists");
    POST("/list/create", { name }, function(e) {
        updateLocalStorage(list, e, "lists");
    });
}

/**
 * Shares a list with a user
 * @param {integer} id of the list
 * @param {string} email of the user to share the list with
 */
function shareList(id, email) {
    POST("/list/share", { id, email });
}

/**
 * Delete a list
 * @param {integer} id of the list
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
        e.forEach(function(elt) { insertLocalStorage(elt, "consumables"); });
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
    var fragmentIndex = listID.indexOf("#");
    if (fragmentIndex !== -1) {
        listID = listID.substring(0, fragmentIndex);
    }
    createCookie("ListID", listID);
    return Number(listID);
}

/**
 * hide all consimables of the page
 */
function hideAllConsumables() {
    var consumables = document.querySelector("tbody");
    while (consumables.firstChild) {
        consumables.removeChild(consumables.firstChild);
    }
}

/**
 * hide a consimable in the list of the page
 * @param {integer} consumableId id of the consumable to remove
 */
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

/**
 * toggle a consumable to be done or not.
 * set the localstorage
 * @param {integer} consumableId id of the consumable to toggle
 */
function toggledoneConsumable(consumableId) {
    c1 = idToElt(consumableId, "consumables");
    c2 = c1;
    c2.done = !c1.done;
    updateLocalStorage(c1, c2, "consumables");
}

/**
 * create and return a new Unic IDentifier.
 */
function getNewUID() {
    var ID = Number(localStorage.getItem("UID")) + 1;
    localStorage.setItem("UID", ID);
    return ID;
}

/**
 * create and return a function to be called with a event listener
 * and stores the consumableID inside.
 * @param {integer} consumableId id of the consumable to validate and be stored
 * in the function
 * @returns {function} to be called from the event listener.
 *
 * @todo update the api
 * @body we might add another function to the api
 */
function validateConsumable(consumableId) {
    return function(event) {
        toggledoneConsumable(consumableId);
        fillConsumables();
    }
}

/**
 * add a consumable to the page
 * @param {object} consumable to display
 */
function showConsumable(consumable) {
    if (!"content" in document.createElement("template") ||
        consumable.ID === undefined || consumable.mode === "sample") { return; }
    var template = document.querySelector("#consumable");
    var clone = document.importNode(template.content, true);
    var td = clone.querySelectorAll("td");
    var tr = clone.querySelector("tr");
    tr.id = consumable.ID
    td[0].addEventListener("click", validateConsumable(consumable.ID));
    if (consumable.done) {
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

/**
 * toggle loading gif depending on the status.
 * uses loadingGifTimeoutID to keep track of timeout IDs.
 * @param {boolean} status display status of the loading gif
 * @param {integer} timeout timeout before showing the gif
 * @param {string} id id the loading gif element
 */
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

/**
 * add a list of consumables to local storage and fill it to the page
 * @param {object} consumables list of consumables to add to local storage
 */
function addConsumablesFromList(consumables) {
    if (consumables.error !== undefined) { return; }
    consumables.forEach(function(elt) { insertLocalStorage(elt, "consumables") });
    fillConsumables();
}

/**
 * fetch all consumables from all lists starting with the current
 */
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
    if (list.id === 0) { return; };
    var json = JSON.parse(localStorage.getItem("consumables"));
    if (json === null) { json = []; }
    showLoadingGif(true);
    var done = [];
    json.sort(function(a, b) {
        return a.ID - b.ID;
    });
    hideAllConsumables();
    json.forEach(function(element) {
        if (element.list_id === list.id) {
            if (element.done) {
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

/**
 * set cookies with the user coordinates
 */
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

/**
 * make te window location to the last list seen
 * @param {string} fragment page fragment to target
 */
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

/**
 * return a new promise to notify the user
 * @param {object} notif notification to display to the user
 * @param {integer} notif.when notification timeout before showing the
 * notification.
 * @return {promise} a new promise that will make a notification.
 */
function promiseNofity(notif) {
    new Promise(function(resolve, reject) {
        setTimeout(function(e) { resolve() }, notif.when);
    }).then(function(event) {
        notify(notif.title, notif.body, notif.onClickURL, notif.icon);
        removeLocalStorage(notif, "notifications");
    });
}

/**
 * register a notification to send to the user in time. it uses the local
 * storage to store new notifications in case of window changing.
 * @param {integer} when notification timeout before showing the notification.
 * @param {string} title title of the notification.
 * @param {string} body body of the notification.
 * @param {string} onclickURL url to redirect to when the notification is
 * clicked.
 * @param {string} icon icon of the notification.
 */
function registerNewNotification(when, title, body, onClickURL = "#",
    icon = "/favicon.ico") {
    var ID = getNewUID();
    notification = { when, title, body, onClickURL, icon, ID };
    promiseNofity(notification);
    insertLocalStorage(notification, "notifications");
}

/**
 * register all previouslly registered notifications from local storage
 */
function registerNotifications() {
    var notifications = JSON.parse(localStorage.getItem("notifications"));
    if (notifications === null) { return; }
    notifications.forEach(promiseNofity);
}

/**
 * send a toast to the user
 * @param {object} t toast to send.
 * @param {integer} t.ID unique identifier.
 * @param {integer} t.pagecount number of page to wait before toasting the user.
 */
function toastMe(t) {
    var tt = idToElt(t.ID, "toasts");
    oldToast = tt;
    if (toast.pagecount-- === 0) {
        toast(tt);
        updateLocalStorage(oldToast, tt, "toast");
    } else {
        removeLocalStorage(tt, "toasts");
    }
}

/**
 * register a toast
 * @param {object} toast toast to send.
 * @param {integer} pagecount number of page to wait before toasting the user.
 */
function registerNewToast(pagecount, toast) {
    toast.pagecount = pagecount;
    toast.ID = getNewUID();
    insertLocalStorage(toast, "toasts");
}

/**
 * call toastMe on each registered toast.
 */
function checkToast() {
    var toasts = JSON.parse(localStorage.getItem("toasts"));
    if (toasts === null) { return; }
    toasts.forEach(toastMe);
}

/**
 * initialise the page load, register notifications and toasts
 */
function init() {
    if (location.search !== "") {
        log("location.search=" + location.search.substring(1));
    }
    registerNotifications()
    checkToast();
}

/**
 * checks if the user is using a mobile or a tablet
 * @return {boolean}
 */
function mobileAndTabletcheck() {
    var check = false;
    (function(a) { if (/(android|bb\d+|meego).+mobile|avantgo|bada\/|blackberry|blazer|compal|elaine|fennec|hiptop|iemobile|ip(hone|od)|iris|kindle|lge |maemo|midp|mmp|mobile.+firefox|netfront|opera m(ob|in)i|palm( os)?|phone|p(ixi|re)\/|plucker|pocket|psp|series(4|6)0|symbian|treo|up\.(browser|link)|vodafone|wap|windows ce|xda|xiino|android|ipad|playbook|silk/i.test(a) || /1207|6310|6590|3gso|4thp|50[1-6]i|770s|802s|a wa|abac|ac(er|oo|s\-)|ai(ko|rn)|al(av|ca|co)|amoi|an(ex|ny|yw)|aptu|ar(ch|go)|as(te|us)|attw|au(di|\-m|r |s )|avan|be(ck|ll|nq)|bi(lb|rd)|bl(ac|az)|br(e|v)w|bumb|bw\-(n|u)|c55\/|capi|ccwa|cdm\-|cell|chtm|cldc|cmd\-|co(mp|nd)|craw|da(it|ll|ng)|dbte|dc\-s|devi|dica|dmob|do(c|p)o|ds(12|\-d)|el(49|ai)|em(l2|ul)|er(ic|k0)|esl8|ez([4-7]0|os|wa|ze)|fetc|fly(\-|_)|g1 u|g560|gene|gf\-5|g\-mo|go(\.w|od)|gr(ad|un)|haie|hcit|hd\-(m|p|t)|hei\-|hi(pt|ta)|hp( i|ip)|hs\-c|ht(c(\-| |_|a|g|p|s|t)|tp)|hu(aw|tc)|i\-(20|go|ma)|i230|iac( |\-|\/)|ibro|idea|ig01|ikom|im1k|inno|ipaq|iris|ja(t|v)a|jbro|jemu|jigs|kddi|keji|kgt( |\/)|klon|kpt |kwc\-|kyo(c|k)|le(no|xi)|lg( g|\/(k|l|u)|50|54|\-[a-w])|libw|lynx|m1\-w|m3ga|m50\/|ma(te|ui|xo)|mc(01|21|ca)|m\-cr|me(rc|ri)|mi(o8|oa|ts)|mmef|mo(01|02|bi|de|do|t(\-| |o|v)|zz)|mt(50|p1|v )|mwbp|mywa|n10[0-2]|n20[2-3]|n30(0|2)|n50(0|2|5)|n7(0(0|1)|10)|ne((c|m)\-|on|tf|wf|wg|wt)|nok(6|i)|nzph|o2im|op(ti|wv)|oran|owg1|p800|pan(a|d|t)|pdxg|pg(13|\-([1-8]|c))|phil|pire|pl(ay|uc)|pn\-2|po(ck|rt|se)|prox|psio|pt\-g|qa\-a|qc(07|12|21|32|60|\-[2-7]|i\-)|qtek|r380|r600|raks|rim9|ro(ve|zo)|s55\/|sa(ge|ma|mm|ms|ny|va)|sc(01|h\-|oo|p\-)|sdk\/|se(c(\-|0|1)|47|mc|nd|ri)|sgh\-|shar|sie(\-|m)|sk\-0|sl(45|id)|sm(al|ar|b3|it|t5)|so(ft|ny)|sp(01|h\-|v\-|v )|sy(01|mb)|t2(18|50)|t6(00|10|18)|ta(gt|lk)|tcl\-|tdg\-|tel(i|m)|tim\-|t\-mo|to(pl|sh)|ts(70|m\-|m3|m5)|tx\-9|up(\.b|g1|si)|utst|v400|v750|veri|vi(rg|te)|vk(40|5[0-3]|\-v)|vm40|voda|vulc|vx(52|53|60|61|70|80|81|83|85|98)|w3c(\-| )|webc|whit|wi(g |nc|nw)|wmlb|wonu|x700|yas\-|your|zeto|zte\-/i.test(a.substr(0, 4))) check = true; })(navigator.userAgent || navigator.vendor || window.opera);
    return check;
};