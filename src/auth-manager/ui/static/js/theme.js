var theme;

$(document).ready(function() {
    theme = localStorage.getItem('velocimodel-theme');
    if (theme) {
        if (theme == 'light') {
            $('.dark').addClass('light').removeClass('dark');
        }
    } else {
        theme = 'light'
        localStorage.setItem('velocimodel-theme', theme);
    }
})

function toggleTheme() {
    if (theme == 'light') {
        theme = 'dark'
        $('.light').addClass('dark').removeClass('light');
        if (typeof editor === 'undefined' || editor === null) {
            console.log("Editor not found!")
        } else {
            console.log("Editor found!")
            monaco.editor.setTheme('velocimodelDarkTheme');
        }
    } else {
        theme = 'light'
        $('.dark').addClass('light').removeClass('dark');
        if (typeof editor === 'undefined' || editor === null) {
            console.log("Editor not found!")
        } else {
            console.log("Editor found!")
            monaco.editor.setTheme('velocimodelLightTheme');
        }
    }
    localStorage.setItem('velocimodel-theme', theme);
}
