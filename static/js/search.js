function displaySearchResult(consumable) {
    console.log(consumable);
}

function search(element) {
    return function(e) {
        console.log(e)
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
    console.log("Adding ", value)
    if (value === undefined) { return; }
    var consumables = localStorage.getItem("consumables");
    var json = [];
    if (consumables !== null) {
        json = JSON.parse(consumables);
    }

    var consumable = {
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
    goBackList();
}

window.addEventListener("load", function() {
    if (!isTokenValid(getCookie("Token"))) {
        window.location.replace("/login.html");
    }
    if (location.search !== "") {
        log("location.search=" + location.search.substring(1));
        elems.value = location.search.substring(3);
        document.getElementById("autocomplete-input-label").classList.add("active");
        search(elems)();
    }
    var elems = document.getElementById("autocomplete-input");
    elems.addEventListener("input", search(elems));
    elems.addEventListener("keyup", function(event) {
        if (event.key === "Enter") {
            addPersonnal(elems.value);
        }
    });
    document.getElementById("nav-back").addEventListener("click", goBackList);
});