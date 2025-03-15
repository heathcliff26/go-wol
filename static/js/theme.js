'use strict';

let storedTheme = ""

function getPreferredTheme() {
    if (storedTheme === "") {
        storedTheme = window.matchMedia('(prefers-color-scheme: dark)').matches ? 'dark' : 'light';
    }
    return storedTheme;
}

function setTheme(theme) {
    document.documentElement.setAttribute('data-bs-theme', theme);
}

// Set the theme on first load
setTheme(getPreferredTheme())

// Watch if the theme changes
window.matchMedia('(prefers-color-scheme: dark)').addEventListener('change', event => {
    const newTheme = event.matches ? "dark" : "light";
    if (newTheme !== storedTheme) {
        storedTheme = newTheme
        setTheme(getPreferredTheme())
    }
})
