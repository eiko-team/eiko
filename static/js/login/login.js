window.addEventListener("load", function() {
    if (isTokenValid(getCookie("Token"))) {
        log("redirect from login");
        redirect();
    } else {
        log("welcome login");
    }

    if (location.search !== "") {
        log("location.search=" + location.search.substring(1));
    }

});

function redirect() {
    deleteCookie("pass_score");
    window.location.replace("/");
}

function closeLogin() {
    setStyleByID("login", "display: none;");
}

function closeRegister() {
    deleteCookie("pass_score");
    setStyleByID("register", "display: none;");
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
};

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
    var score = getCookie("pass_score");
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

function displayLoadingGif(display = false, id = "login") {
    if (display) {
        document.getElementById(id + "-button").disabled = true;
        document.getElementById(id + "-loading-gif").style.display = "";
        document.getElementById(id + "-loading-text").style.display = "none";
    } else {
        document.getElementById(id + "-button").disabled = false;
        document.getElementById(id + "-loading-gif").style.display = "none";
        document.getElementById(id + "-loading-text").style.display = "";
    }
}

function login(email, password, remember = true) {
    displayLoadingGif(true, "login");
    POST("/login", { user_email: email, user_password: password }, (e) => {
        if (e.token === undefined) {
            displayLoadingGif(false, "login");
            setStyleByID("error-email", "style: ;");
            return false;
        }
        var user_token = e.token;
        createCookie("token", user_token, remember ? 7 : null);
        log("login");
        redirect();
        return true;
    });
}

function register(email, password, remember = true) {
    displayLoadingGif(true, "register");
    POST("/register", { user_email: email, user_password: password }, (e) => {
        if (e.token === undefined) {
            displayLoadingGif(false, "register");
            setStyleByID("error-email-register", "style: ;");
            return false;
        }
        var user_token = e.token;
        createCookie("token", user_token, remember ? 7 : null);
        log("register");
        redirect();
    });
}

function loginForm() {
    var email = document.forms["login"]["email"].value;
    var password = document.forms["login"]["password"].value;
    var remember = document.forms["login"]["remember"].value;
    login(email, password, remember === "on");
}

function registerForm() {
    var email = document.forms["register"]["email"].value;
    var password1 = document.forms["register"]["password1"].value;
    var password2 = document.forms["register"]["password2"].value;
    var remember = document.forms["register"]["rmb"].value
    if (password1 === password2) {
        register(email, password1, remember === "on");
    }
}