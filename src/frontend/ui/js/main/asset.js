// import material.js
// import theme.js
// import modal.js

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
        url: "/script/api/asset/" + assetID,
        type: "PUT",
        contentType: 'application/json',
        data: JSON.stringify(data),
        success: function(response) {
            $("#spinner").css("display", "none")
            $("#page-darken").css("opacity", "0")
        },
        error: function(response) {
            console.log(response)
            $("#log-container").html(response.responseJSON['output'])
            $("#spinner").css("display", "none")
            $("#page-darken").css("opacity", "0")
            openModal('error-modal')
        }
    });
}
