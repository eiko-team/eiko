setTheme(getCookie("theme"));

function swapTheme() {
    var theme = "light";
    if (getCookie("theme") !== "dark") {
        theme = "dark";
    }
    setTheme(theme);
}

function setTheme(theme = "light") {
    createCookie("theme", theme);
    document.documentElement.classList.add("color-theme-in-transition");
    document.documentElement.setAttribute("data-theme", theme);
    window.setTimeout(function() {
        document.documentElement.classList.remove("color-theme-in-transition");
    }, 1000);
}