
// import material.js
// import theme.js

function getReleases() {
    $.ajax({
        url: "/script/api/release",
        type: "GET",
        success: function (result) {
            render(result);
        }
    });
}

function render(data) {
    var prefix = $("#search").val();
    if (prefix != "") {
        for (var idx = data.length - 1; idx >= 0; idx--) {
            if (!data[idx].name.startsWith(prefix)) {
                data.splice(idx, 1)
            }
        }
    }
    var table = document.getElementById('releases-table');
    var tableHTMLString = '<tr><th>Release Name</th><th>Last Updated</th><th>Language</th><th>Tags</th><th></th></tr>' +
        data.map(function (release) {
            return '<tr>' +
                '<td>' + release.name + '</td>' +
                '<td>' + release.updated + '</td>' +
                '<td>' + release.language + '</td>' +
                '<td>' + release.tags.map(function (tag) {
                    return '<span class="w3-tag w3-round velocimodel-blue release-tag">' + tag + '</span>'
                }).join('') +
                '</td>' +
                '<td>' +
                '<a href="/ui/release/' + release.id + '"><i class="fa-solid fa-link"></i></a>' +
                '</td>' +
                '</tr>'
        }).join('');

    table.innerHTML = tableHTMLString;
}
