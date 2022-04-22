
// import material.js
// import theme.js

function getModels() {
    $.ajax({
        url: "/script/api/model",
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
    var table = document.getElementById('models-table');
    var tableHTMLString = '<tr><th class="table-title w3-medium velocimodel-text-blue"><span class="table-title-text">Model Name</span></th><th class="table-title w3-medium velocimodel-text-blue"><span class="table-title-text">Last Updated</span></th><th class="table-title w3-medium velocimodel-text-blue"><span class="table-title-text">Language</span></th><th class="table-title w3-medium velocimodel-text-blue"><span class="table-title-text">Tags</span></th><th><div class="w3-round w3-button velocimodel-green">Create New</div></th></tr>' +
        data.map(function (model) {
            return '<tr>' +
                '<td>' + model.name + '</td>' +
                '<td>' + model.updated + '</td>' +
                '<td>' + model.language + '</td>' +
                '<td>' + model.tags.map(function (tag) {
                    return '<span class="w3-tag w3-round velocimodel-green asset-tag">' + tag + '</span>'
                }).join('') +
                '</td>' +
                '<td class="table-link-cell">' +
                '<a href="/ui/model/' + model.id + '" class="table-link-link w3-right-align dark theme-text" style="float:right;margin-right:16px;"><i class="fa-solid fa-link"></i></a>' +
                '</td>' +
                '</tr>'
        }).join('');

    table.innerHTML = tableHTMLString;
}
