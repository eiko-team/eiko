if ("serviceWorker" in navigator) {
    if (!navigator.serviceWorker.controller) {
        navigator.serviceWorker.register("/eiko-sw.js");
    }
}

window.addEventListener("load", function() {
    if (!isTokenValid(getCookie("Token"))) {
        window.location.replace("/login.html");
    }
    loadList();
    if (location.search !== "") {
        log("location.search=" + location.search.substring(1));
    }
});


window.addEventListener("keydown", function(e) {
    if ((e.key === "Escape" || e.key === "Esc" || e.keyCode === 27) &&
        (e.target.nodeName === "BODY")) {
        closeNav();
        e.preventDefault();
        return false;
    }
}, true);

var dropdown = document.getElementsByClassName("dropdown-btn");
for (var i = 0; i < dropdown.length; i++) {
    dropdown[i].addEventListener("click", function() {
        this.classList.toggle("active");
        var dropdownContent = this.nextElementSibling;
        if (dropdownContent.style.display === "block") {
            this.getElementsByClassName("material-icons")[1].innerText = "keyboard_arrow_down";
            dropdownContent.style.display = "none";
        } else {
            this.getElementsByClassName("material-icons")[1].innerText = "keyboard_arrow_up";
            dropdownContent.style.display = "block";
        }
    });
}

if ("content" in document.createElement("template")) {
    var tbody = document.querySelector("tbody");
    var template = document.querySelector("#consumable");
    var clone = document.importNode(template.content, true);
    var td = clone.querySelectorAll("td");
    td[0].textContent = "1235646565";
    td[1].textContent = "Stuff";
    tbody.appendChild(clone);

    var clone2 = document.importNode(template.content, true);
    td = clone2.querySelectorAll("td");
    td[0].textContent = "0384928528";
    td[1].textContent = "Acme Kidney Beans";
    tbody.appendChild(clone2);
} else {
    // no templates available
}

document.getElementById("logout-button").addEventListener("click", logout);
document.getElementById("theme-selector").addEventListener("click", swapTheme);
document.getElementById("open-nav").addEventListener("click", openNav);
document.getElementById("nav-closebtn").addEventListener("click", closeNav);