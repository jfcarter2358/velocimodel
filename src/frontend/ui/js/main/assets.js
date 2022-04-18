
// import material.js
// import theme.js

function getAssets() {
    $.ajax({
        url: "/script/api/asset",
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
    var table = document.getElementById('assets-table');
    var tableHTMLString = '<tr><th>Asset Name</th><th>Last Updated</th><th>Language</th><th>Tags</th><th></th></tr>' +
        data.map(function (asset) {
            return '<tr>' +
                '<td>' + asset.name + '</td>' +
                '<td>' + asset.updated + '</td>' +
                '<td>' + asset.language + '</td>' +
                '<td>' + asset.tags.map(function (tag) {
                    return '<span class="w3-tag w3-round velocimodel-blue asset-tag">' + tag + '</span>'
                }).join('') +
                '</td>' +
                '<td>' +
                '<a href="/ui/asset/' + asset.id + '"><i class="fa-solid fa-link"></i></a>' +
                '</td>' +
                '</tr>'
        }).join('');

    table.innerHTML = tableHTMLString;
}
