// import material.js
// import theme.js
// import modal.js

var tagify

$(document).ready(
    function() {
        var input = document.getElementById('tag-input'),
        tagify = new Tagify(input, {
            id: 'model_tags',
        })
        var coll = document.getElementsByClassName("collapsible");
        var i;

        for (i = 0; i < coll.length; i++) {
            coll[i].addEventListener("click", function() {
                this.classList.toggle("active");
                var content = this.nextElementSibling;
                if (content.style.maxHeight){
                content.style.maxHeight = null;
                } else {
                content.style.maxHeight = content.scrollHeight + "px";
                } 
            });
        }
    }
)

function saveModel() {
    tagData = []
    if (document.getElementById('tag-input').value != "") {
        tagData = JSON.parse(document.getElementById('tag-input').value)
    }

    parts = window.location.href.split('/')
    modelID = parts[parts.length - 1]

    modelName = $("#model-name").val();
    tags = []
    for (var i = 0; i < tagData.length; i++) {
        tags.push(tagData[i]["value"])
    }
    assets = []
    $(".asset-id").each((index, elem) => {
        assets.push(elem.innerText);
    });
    
    data = {
        "name": modelName,
        "tags": tags,
        "assets": assets
    }

    $("#spinner").css("display", "block")
    $("#page-darken").css("opacity", "1")

    $.ajax({
        url: "/script/api/model/" + modelID,
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

function createSnapshot() {
    parts = window.location.href.split('/')
    modelID = parts[parts.length - 1]

    $("#spinner").css("display", "block")
    $("#page-darken").css("opacity", "1")

    $.ajax({
        url: "/script/api/model/" + modelID + "/snapshot",
        type: "POST",
        success: function(response) {
            $("#spinner").css("display", "none")
            $("#page-darken").css("opacity", "0")
            window.location.reload();
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

function openExistingModal() {
    closeModal('add-asset-modal')
    openModal('add-existing-asset-modal')
}

function openFileModal() {
    closeModal('add-asset-modal')
    openModal('add-file-asset-modal')
}

function openGitModal() {
    closeModal('add-asset-modal')
    openModal('add-git-asset-modal')
}

function addExistingAsset() {
    parts = window.location.href.split('/')
    modelID = parts[parts.length - 1]
    assetID = $("#existing-asset-id").val()
    
    data = {
        "model": modelID,
        "asset": assetID
    }

    $("#spinner").css("display", "block")
    $("#page-darken").css("opacity", "1")

    $.ajax({
        url: "/script/api/model/asset",
        type: "POST",
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

function addGitAsset() {
    // repo
    // branch
    // public
}

function addFileAsset() {
    $('#file-form')
    .ajaxForm({
        url : '/script/api/asset/file',
        type: "POST",
        success : function (response) {
            $("#spinner").css("display", "none")
            $("#page-darken").css("opacity", "0")
            console.log(response)
        },
        error: function(response) {
            console.log(response)
            $("#log-container").html(response.responseJSON['output'])
            $("#spinner").css("display", "none")
            $("#page-darken").css("opacity", "0")
            openModal('error-modal')
        }
    })
;
}