function displaySearchResult(consumable) {}
function clearSearchResult(consumable) {}

function search(element) {
    return function(e) {
        showLoadingGif(true)
        getTargetPosition();
        POST("/consumable/get", {
            query: element.value,
            latitude: getCookie("posLat"),
            longitude: getCookie("posLon"),
            accuracy: getCookie("posAcc")
        }, function(e) {
            e.forEach(displaySearchResult);
            showLoadingGif(false)
        }, function(e) {
            showLoadingGif(false)
        });
    }
}

function addPersonnal(value) {
    if (value === undefined) { return; }
    var consumables = localStorage.getItem("consumables");
    var json = [];
    if (consumables !== null) {
        json = JSON.parse(consumables);
    }
    var uuid = getNewUID();
    var consumable = {
        uuid: uuid,
        ID: uuid,
        ListID: getCookie("ListID"),
        Name: value,
        Done: false,
        Erased: false,
        Mode: "personnal",
    }

    json.push(consumable)
    localStorage.setItem("consumables", JSON.stringify(json));
    POST("/list/add/personnal", consumable, function(event) {
        var json = JSON.parse(localStorage.getItem("consumables"));
        if (json === null) { return; };
        json.forEach(function(element) {
            if (element.ID === consumable.ID) {
                element.ID = event.ID
            }
        });
        localStorage.setItem("consumables", JSON.stringify(json));
    })
    window.history.back();
}

function fillAutoComplete() {
    GET("/json/autocomplete_data.json", function(data) {
        M.Autocomplete.init(document.querySelectorAll('.autocomplete'), { data });
    });
}

window.addEventListener("DOMContentLoaded", function() {
    if (!isTokenValid(getCookie("Token"))) {
        window.location.replace("/login.html");
    }
    if (location.search !== "") {
        log("location.search=" + location.search.substring(1));
        elems.value = location.search.substring(3);
        document.getElementById("search-input-label").classList.add("active");
        search(elems)();
    }
    var addBtn = document.getElementById("add-item");
    addBtn.addEventListener("click", function(e) {
        addPersonnal(elems.value);
    });
    fillAutoComplete();
    var elems = document.getElementById("search-input");
    elems.addEventListener("input", search(elems));
    elems.focus();
    elems.select();
    elems.addEventListener("keyup", function(event) {
        if (event.key === "Enter") {
            addPersonnal(elems.value);
        }
    });
    document.getElementById("nav-back").addEventListener("click", function(e) {
        window.history.back();
    });
});