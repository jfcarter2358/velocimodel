
// import material.js
// import theme.js
// import modal.js

var assets

function getAssets() {
    $.ajax({
        url: "/api/asset",
        type: "GET",
        success: function (result) {
            assets = result
        }
    });
}

function render() {
    var prefix = $("#search").val();
    var tempAssets = JSON.parse(JSON.stringify(assets));
    prefix = prefix.toLowerCase();
    if (prefix != "") {
        for (var idx = tempAssets.length - 1; idx >= 0; idx--) {
            if (tempAssets[idx].name.toLowerCase().indexOf(prefix) == -1) {
                tempAssets.splice(idx, 1)
            }
        }
    }
    var table = document.getElementById('assets-table');
    var tableHTMLString = '<tr><th class="table-title w3-medium velocimodel-text-blue"><span class="table-title-text">Asset Name</span></th><th class="table-title w3-medium velocimodel-text-blue"><span class="table-title-text">Last Updated</span></th><th class="table-title w3-medium velocimodel-text-blue"><span class="table-title-text">Type</span></th><th class="table-title w3-medium velocimodel-text-blue"><span class="table-title-text">Tags</span></th><th><div class="w3-round w3-button velocimodel-green">Create New</div></th></tr>' +
        tempAssets.map(function (asset) {
            return '<tr>' +
                '<td>' + asset.name + '</td>' +
                '<td>' + asset.updated + '</td>' +
                '<td>' + asset.type + '</td>' +
                '<td>' + asset.tags.map(function (tag) {
                    return '<span class="w3-tag w3-round velocimodel-green asset-tag">' + tag + '</span>'
                }).join('') +
                '</td>' +
                '<td class="table-link-cell">' +
                '<a href="/ui/asset/' + asset.id + '" class="table-link-link w3-right-align dark theme-text" style="float:right;margin-right:16px;"><i class="fa-solid fa-link"></i></a>' +
                '</td>' +
                '</tr>'
        }).join('');

    table.innerHTML = tableHTMLString;
}

$(document).ready(
    function() {
        getAssets()
    }
)

function openFileModal() {
    closeModal('add-asset-modal')
    openModal('add-file-asset-modal')
}

function openGitModal() {
    closeModal('add-asset-modal')
    openModal('add-git-asset-modal')
}

function addGitAsset() {
    parts = window.location.href.split('/')
    modelID = parts[parts.length - 1]

    data = {
        "repo": $("#git-asset-repo").val(),
        "branch": $("#git-asset-branch").val(),
        "credential": $("#git-asset-credential").val()
    }

    $("#spinner").css("display", "block")
    $("#page-darken").css("opacity", "1")

    $.ajax({
        url: "/api/asset/git",
        type: "POST",
        contentType: 'application/json',
        data: JSON.stringify(data),
        success: function(response) {
            $("#spinner").css("display", "none")
            $("#page-darken").css("opacity", "0")
            assetID = response["id"]
            closeModal('add-git-asset-modal');
            window.location.assign('/ui/asset/' + assetID);
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

function addFileAsset() {
    parts = window.location.href.split('/')
    modelID = parts[parts.length - 1]

    $('#file-form').ajaxSubmit({
        url : '/api/asset/file',
        type: "POST",
        success : function (response) {
            $("#spinner").css("display", "none")
            $("#page-darken").css("opacity", "0")
            assetID = response["id"]
            closeModal('add-file-asset-modal');
            window.location.assign('/ui/asset/' + assetID);
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

$(document).ready(
    function() {
        $('#file-form').ajaxForm()
    }
)
