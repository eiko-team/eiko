if ("serviceWorker" in navigator) {
    if (!navigator.serviceWorker.controller) {
        navigator.serviceWorker.register("/eiko-sw.js", { scope: "./" });
    }
}

function redirect() {
    deleteCookie("pass_score");
    window.location.replace("/");
}

if (isTokenValid(getCookie("Token"))) {
    log("redirect from login");
    redirect();
} else {
    log("login");
    log("welcome login");
}

var login_form = document.getElementById("login");
var register_form = document.getElementById("register");

window.onclick = function(event) {
    if (event.target === login_form) {
        closeLogin();
    }
    if (event.target === register_form) {
        closeRegister();
    }
}

window.addEventListener("keydown", function(e) {
    if ((e.key === "Escape" || e.key === "Esc" || e.keyCode === 27) &&
        (e.target.nodeName === "BODY")) {
        closeLogin();
        closeRegister();
        closeNav();
        e.preventDefault();
        return false;
    }
}, true);

var password = document.getElementById("password1");
password.addEventListener("input", function() {
    var text = document.getElementById("password-strength-text");
    var val = password.value;
    checkPassword(val);
    var score = getCookie("pass_score")
    var meter = document.getElementById("password-strength-meter");

    // Update the password strength meter
    meter.value = score;
    meter.style = "background: " + {
        1: "red",
        2: "yellow",
        3: "orange",
        4: "green",
    } [score] + ";";

    // Update the text indicator
    if (val !== "") {
        text.innerText = "Mot de passe " + {
            0: "mauvais",
            1: "faible",
            2: "moyen",
            3: "bon",
            4: "fort"
        } [score];
    } else {
        text.innerHTML = "";
    }
});

function login(email, password, remember = true) {
    POST("/login", { user_email: email, user_password: password }, (e) => {
        if (e.token === undefined) {
            setStyleByID("error-email", "style: ;");
            return false;
        }
        user_token = e.token;
        createCookie("token", user_token, remember ? 7 : null);
        log("login");
        redirect();
        return true;
    });
}

function register(email, password, remember = true) {
    POST("/register", { user_email: email, user_password: password }, (e) => {
        if (e.token === undefined) {
            setStyleByID("error-email-register", "style: ;");
            return false;
        }
        user_token = e.token;
        createCookie("token", user_token, remember ? 7 : null);
        log("register");
        redirect();
    });
}

function loginForm(email, password, remember) {
    console.log("logging")
    login(email, password, remember === "on");
}

function registerForm(email, password1, password2, remember) {
    if (password1 === password2) {
        register(email, password1, remember === "on");
    }
}