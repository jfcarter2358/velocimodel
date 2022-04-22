
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
    var tableHTMLString = '<tr><th class="table-title w3-medium velocimodel-text-blue">Release Name</th><th class="table-title w3-medium velocimodel-text-blue">Last Updated</th><th class="table-title w3-medium velocimodel-text-blue">Language</th><th class="table-title w3-medium velocimodel-text-blue">Tags</th><th><div class="w3-round w3-button velocimodel-green">Create New</div></th></tr>' +
        data.map(function (release) {
            return '<tr>' +
                '<td>' + release.name + '</td>' +
                '<td>' + release.updated + '</td>' +
                '<td>' + release.language + '</td>' +
                '<td>' + release.tags.map(function (tag) {
                    return '<span class="w3-tag w3-round velocimodel-green asset-tag">' + tag + '</span>'
                }).join('') +
                '</td>' +
                '<td class="table-link-cell">' +
                '<a href="/ui/release/' + release.id + '" class="table-link-link w3-right-align dark theme-text" style="float:right;margin-right:16px;"><i class="fa-solid fa-link"></i></a>' +
                '</td>' +
                '</tr>'
        }).join('');

    table.innerHTML = tableHTMLString;
}
