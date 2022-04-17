
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
    var tableHTMLString = '<tr><th>Model Name</th><th>Last Updated</th><th>Language</th><th>Tags</th><th></th></tr>' +
        data.map(function (model) {
            return '<tr>' +
                '<td>' + model.name + '</td>' +
                '<td>' + model.updated + '</td>' +
                '<td>' + model.language + '</td>' +
                '<td>' + model.tags.map(function (tag) {
                    return '<span class="w3-tag w3-round velocimodel-blue model-tag">' + tag + '</span>'
                }).join('') +
                '</td>' +
                '<td>' +
                '<a href="/ui/model/' + model.id + '"><i class="fa-solid fa-link"></i></a>' +
                '</td>' +
                '</tr>'
        }).join('');

    table.innerHTML = tableHTMLString;
}
