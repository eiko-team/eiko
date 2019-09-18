var user_token = ""

if ("serviceWorker" in navigator) {
    if (!navigator.serviceWorker.controller) {
        navigator.serviceWorker.register("/eiko-sw.js", { scope: "./" });
    }
}

if (!isTokenValid()) {
    setStyleByClass("nav-login", "");
    setStyleByClass("nav-menu", "display: none;");
    setStyleByClass("nav-nav", "display: none;");
}

var login_form = document.getElementById("login");
var register_form = document.getElementById("register");

window.onclick = function(event) {
    if (event.target == login_form) {
        login_form.style.display = "none";
    }
    if (event.target == register_form) {
        register_form.style.display = "none";
    }
}

window.addEventListener('keydown', function(e) {
    if ((e.key == 'Escape' || e.key == 'Esc' || e.keyCode == 27) && (e.target.nodeName == 'BODY')) {
        setStyleByID("login", "display: none;")
        setStyleByID("register", "display: none;")
        closeNav();
        e.preventDefault();
        return false;
    }
}, true);

var password = document.getElementById('password1');
password.addEventListener('input', function() {
    var text = document.getElementById('password-strength-text');
    var val = password.value;
    var result = checkPassword(val);
    var meter = document.getElementById('password-strength-meter');

    // Update the password strength meter
    meter.value = result.score;
    meter.style = "background: " + {
        1: "red",
        2: "yellow",
        3: "orange",
        4: "green",
    } [result.score] + ";"

    // Update the text indicator
    if (val !== "") {
        text.innerHTML = {
            0: "Pire",
            1: "Mauvais",
            2: "Faible",
            3: "Bon",
            4: "Fort"
        } [result.score] + " Mot De Passe";
    } else {
        text.innerHTML = "";
    }
});

function openNav() {
    document.getElementById("mySidenav").style.width = "250px";
}

function closeNav() {
    document.getElementById("mySidenav").style.width = "0";
}

function login(email, password, remember = true) {
    return POST("/login", { user_email: email, user_password: password }, (e) => {
        if (e.token === undefined) {
            setStyleByID("error-email", "style: ;")
            return false
        }
        user_token = e.token
        createCookie("token", user_token, remember ? 7 : null)
        log("login")
        window.location.reload(false);
        return true
    })
}

function register(email, password, remember = true) {
    POST("/register", { user_email: email, user_password: password }, (e) => {
        if (e.token === undefined) {
            setStyleByID("error-email-register", "style: ;")
            return false
        }
        user_token = e.token
        createCookie("token", user_token, remember ? 7 : null)
        log("register")
        window.location.reload(false);
    })
}

function logout() {
    log("logout")
    deleteCookie("token")
    window.location.reload(false);
}

function loginForm(email, password, remember) {
    login(email, password, remember === 'on')
}

function registerForm(email, password1, password2, remember) {
    if (password1 === password2)
        register(email, password1, remember === 'on')
}