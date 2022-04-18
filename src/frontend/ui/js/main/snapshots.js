
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
    var tableHTMLString = '<tr><th>Snapshot Name</th><th>Last Updated</th><th>Language</th><th>Tags</th><th></th></tr>' +
        data.map(function (snapshot) {
            return '<tr>' +
                '<td>' + snapshot.name + '</td>' +
                '<td>' + snapshot.updated + '</td>' +
                '<td>' + snapshot.language + '</td>' +
                '<td>' + snapshot.tags.map(function (tag) {
                    return '<span class="w3-tag w3-round velocimodel-blue snapshot-tag">' + tag + '</span>'
                }).join('') +
                '</td>' +
                '<td>' +
                '<a href="/ui/snapshot/' + snapshot.id + '"><i class="fa-solid fa-link"></i></a>' +
                '</td>' +
                '</tr>'
        }).join('');

    table.innerHTML = tableHTMLString;
}
