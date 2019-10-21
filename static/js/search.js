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
    var uuid = getNewUID();
    var consumable = {
        ID: uuid,
        list_id: Number(getCookie("ListID")),
        name: value,
        done: false,
        erased: false,
        mode: "personnal",
    }

    addConsumableLocalStorage(consumable);
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
    var elems = document.getElementById("search-input");
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