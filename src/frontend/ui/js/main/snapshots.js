
// import material.js
// import theme.js

function getSnapshots() {
    $.ajax({
        url: "/script/api/snapshot",
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
    var table = document.getElementById('snapshots-table');
var tableHTMLString = '<tr><th class="table-title w3-medium velocimodel-text-blue">Snapshot Name</th><th class="table-title w3-medium velocimodel-text-blue">Last Updated</th><th class="table-title w3-medium velocimodel-text-blue">Language</th><th class="table-title w3-medium velocimodel-text-blue">Tags</th><th><div class="w3-round w3-button velocimodel-green">Create New</div></th></tr>' +
        data.map(function (snapshot) {
            return '<tr>' +
                '<td>' + snapshot.name + '</td>' +
                '<td>' + snapshot.updated + '</td>' +
                '<td>' + snapshot.language + '</td>' +
                '<td>' + snapshot.tags.map(function (tag) {
                    return '<span class="w3-tag w3-round velocimodel-green asset-tag">' + tag + '</span>'
                }).join('') +
                '</td>' +
                '<td class="table-link-cell">' +
                '<a href="/ui/snapshot/' + snapshot.id + '" class="table-link-link w3-right-align dark theme-text" style="float:right;margin-right:16px;"><i class="fa-solid fa-link"></i></a>' +
                '</td>' +
                '</tr>'
        }).join('');

    table.innerHTML = tableHTMLString;
}
