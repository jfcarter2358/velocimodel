
// import material.js
// import theme.js
// import modal.js
// import user_menu.js
// import status.js
// import data.js

var models

function getModels() {
    $.ajax({
        url: "/api/model",
        type: "GET",
        success: function (result) {
            models = result
        }
    });
}

function render() {
    var prefix = $("#search").val();
    var tempModels = JSON.parse(JSON.stringify(models));
    prefix = prefix.toLowerCase();
    if (prefix != "") {
        for (var idx = tempModels.length - 1; idx >= 0; idx--) {
            if (tempModels[idx].name.toLowerCase().indexOf(prefix) == -1) {
                tempModels.splice(idx, 1)
            }
        }
    }
    var table = document.getElementById('models-table');
    var tableHTMLString = '<tr><th class="table-title w3-medium velocimodel-text-blue"><span class="table-title-text">Model Name</span></th><th class="table-title w3-medium velocimodel-text-blue"><span class="table-title-text">Last Updated</span></th><th class="table-title w3-medium velocimodel-text-blue"><span class="table-title-text">Language</span></th><th class="table-title w3-medium velocimodel-text-blue"><span class="table-title-text">Tags</span></th><th><div class="w3-round w3-button velocimodel-green">Create New</div></th></tr>' +
        tempModels.map(function (model) {
            return '<tr>' +
                '<td>' + model.name + '</td>' +
                '<td>' + model.updated + '</td>' +
                '<td>' + model.language + '</td>' +
                '<td>' + model.tags.map(function (tag) {
                    return '<span class="w3-tag w3-round velocimodel-green asset-tag">' + tag + '</span>'
                }).join('') +
                '</td>' +
                '<td class="table-link-cell">' +
                '<a href="/ui/model/' + model.id + '" class="table-link-link w3-right-align light theme-text" style="float:right;margin-right:16px;"><i class="fa-solid fa-link"></i></a>' +
                '</td>' +
                '</tr>'
        }).join('');

    table.innerHTML = tableHTMLString;
}

$(document).ready(
    function() {
        getModels()
    }
)

function createModel() {
    data = {
        "id": "",
        "name": $("#models-create-name").val(),
        "created": "",
        "updated": "",
        "type": "raw",
        "tags": [],
        "metadata": {},
        "assets": [],
        "snapshots": [],
        "releases": [],
        "language": $("#models-create-language").val()
    }

    $("#spinner").css("display", "block")
    $("#page-darken").css("opacity", "1")

    $.ajax({
        url: "/api/model",
        type: "POST",
        contentType: 'application/json',
        data: JSON.stringify(data),
        success: function(response) {
            $("#spinner").css("display", "none")
            $("#page-darken").css("opacity", "0")
            modelID = response['id']
            closeModal('models-create-modal');
            window.location.assign('/ui/model/' + modelID);
        },
        error: function(response) {
            console.log(response)
            $("#log-container").text(response.responseJSON['error'])
            $("#spinner").css("display", "none")
            $("#page-darken").css("opacity", "0")
            openModal('error-modal')
        }
    });
}
