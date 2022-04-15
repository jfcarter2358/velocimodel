var theme;

$(document).ready(function() {
    theme = localStorage.getItem('moc-builder-theme');
    if (theme) {
        if (theme == 'light') {
            $('.dark').addClass('light').removeClass('dark');
            // $("#theme-button").addClass('fa-sun').removeClass('fa-moon')
        }
    } else {
        theme = 'light'
        localStorage.setItem('moc-builder-theme', theme);
    }
})

function toggleTheme() {
    if (theme == 'light') {
        theme = 'dark'
        // $("#theme-button").addClass('fa-moon').removeClass('fa-sun')
        $('.light').addClass('dark').removeClass('light');
    } else {
        theme = 'light'
        // $("#theme-button").addClass('fa-sun').removeClass('fa-moon')
        $('.dark').addClass('light').removeClass('dark');
    }
    localStorage.setItem('moc-builder-theme', theme);
}