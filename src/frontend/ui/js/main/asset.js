// import material.js
// import theme.js
// import modal.js
// import user_menu.js
// import status.js
// import data.js

var tagify

$(document).ready(
    function() {
        var input = document.getElementById('tag-input'),
        tagify = new Tagify(input, {
            id: 'asset_tags',
        })
    }
)

function saveAsset() {
    console.log(document.getElementById('tag-input').value)
    tagData = []
    if (document.getElementById('tag-input').value != "") {
        tagData = JSON.parse(document.getElementById('tag-input').value)
    }

    parts = window.location.href.split('/')
    basePath = ''
if (parts[3] != 'ui') {
        basePath = "/" + parts[3]
    }
    assetID = parts[parts.length - 1]

    assetName = $("#asset-name").val();
    tags = []
    for (var i = 0; i < tagData.length; i++) {
        tags.push(tagData[i]["value"])
    }

    data = {
        "name": assetName,
        "tags": tags
    }

    $("#spinner").css("display", "block")
    $("#page-darken").css("opacity", "1")

    $.ajax({
        url: basePath + "/api/asset/" + assetID,
        type: "PUT",
        contentType: 'application/json',
        data: JSON.stringify(data),
        success: function(response) {
            $("#spinner").css("display", "none")
            $("#page-darken").css("opacity", "0")
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

function downloadFileAsset(url) {
    const a = document.createElement('a')
    a.href = url
    a.download = url.split('/').pop()
    document.body.appendChild(a)
    a.click()
    document.body.removeChild(a)
}

function syncGitAsset(assetID) {
    parts = window.location.href.split('/')
    basePath = ''
if (parts[3] != 'ui') {
        basePath = "/" + parts[3]
    }

    $("#spinner").css("display", "block")
    $("#page-darken").css("opacity", "1")

    $.ajax({
        url: basePath + "/api/asset/git/sync/" + assetID,
        type: "POST",
        success: function(response) {
            $("#spinner").css("display", "none")
            $("#page-darken").css("opacity", "0")
            window.location.reload();
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
