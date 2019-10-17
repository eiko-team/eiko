if ("serviceWorker" in navigator) {
    if (!navigator.serviceWorker.controller) {
        navigator.serviceWorker.register("/eiko-sw.js");
    }
}

function showButton() {
    var elems = document.querySelectorAll(".fixed-action-btn");
    var mobile = mobileAndTabletcheck();
    var addBtn = document.getElementById("add-item-floating")
    if (mobile) {
        addBtn.style.display = "";
    } else {
        document.getElementById("add-item").href = addBtn.href;
    }
    var instances = M.FloatingActionButton.init(elems, {
        direction: "top",
        hoverEnabled: !mobile,
        toolbarEnabled: false
    });
}

window.addEventListener("DOMContentLoaded", function() {
    if (!isTokenValid(getCookie("Token"))) {
        window.location.replace("/login.html");
    }
    init();
    loadLists();
    fillConsumables(true);
    if (location.search !== "") {
        log("location.search=" + location.search.substring(1));
    }
    showButton();
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

document.getElementById("logout-button").addEventListener("click", logout);
document.getElementById("theme-selector").addEventListener("click", swapTheme);
document.getElementById("open-nav").addEventListener("click", openNav);
document.getElementById("nav-closebtn").addEventListener("click", closeNav);