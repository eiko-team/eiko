function displaySearchResult(consumable) {}

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

var UUID = 1;

function addPersonnal(value) {
    if (value === undefined) { return; }
    var consumables = localStorage.getItem("consumables");
    var json = [];
    if (consumables !== null) {
        json = JSON.parse(consumables);
    }

    var consumable = {
        uuid: UUID,
        ID: UUID++,
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

var autocompleteData = {
    "Carrotes": null,
    "Dentifrice": null,
    "Pain": null,
    "Pastèque": null,
    "Patates": null,
    "Poires": null,
    "Poivons": null,
    "Pommes": null
}

window.addEventListener("load", function() {
    if (!isTokenValid(getCookie("Token"))) {
        window.location.replace("/login.html");
    }
    if (location.search !== "") {
        log("location.search=" + location.search.substring(1));
        elems.value = location.search.substring(3);
        document.getElementById("search-input-label").classList.add("active");
        search(elems)();
    }
    var elems = document.getElementById("search-input");
    elems.addEventListener("input", search(elems));
    elems.focus();
    elems.select();
    elems.addEventListener("keyup", function(event) {
        if (event.key === "Enter") {
            addPersonnal(elems.value);
        }
    });
    var addBtn = document.getElementById("add-item");
    addBtn.addEventListener("click", function(e) {
        addPersonnal(elems.value);
    });
    var autocompleteElt = document.querySelectorAll('.autocomplete');
    var autocomplete = M.Autocomplete.init(autocompleteElt, {
        data: autocompleteData
    });
    document.getElementById("nav-back").addEventListener("click",
        window.history.back);
});