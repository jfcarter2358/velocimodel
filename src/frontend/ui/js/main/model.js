// import material.js
// import theme.js
// import modal.js

var tagify

document.addEventListener("DOMContentLoaded", function (event) {
    var scrollpos = sessionStorage.getItem('velocimodelscrollpos');
    if (scrollpos) {
        window.scrollTo(0, scrollpos);
        sessionStorage.removeItem('velocimodelscrollpos');
    }
});

window.addEventListener("beforeunload", function (e) {
    sessionStorage.setItem('velocimodelscrollpos', window.scrollY);
});

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
        url: "/api/model/" + modelID,
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

function createSnapshot() {
    parts = window.location.href.split('/')
    modelID = parts[parts.length - 1]

    $("#spinner").css("display", "block")
    $("#page-darken").css("opacity", "1")

    $.ajax({
        url: "/api/model/" + modelID + "/snapshot",
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

function addExistingAsset(assetID) {
    parts = window.location.href.split('/')
    modelID = parts[parts.length - 1]
    
    data = {
        "model": modelID,
        "asset": assetID
    }

    $("#spinner").css("display", "block")
    $("#page-darken").css("opacity", "1")

    $.ajax({
        url: "/api/model/asset",
        type: "POST",
        contentType: 'application/json',
        data: JSON.stringify(data),
        success: function(response) {
            $("#spinner").css("display", "none")
            $("#page-darken").css("opacity", "0")
            closeModal('add-existing-asset-modal');
            window.location.reload();
        },
        error: function(response) {
            console.log(response)
            $("#log-container").text(response.responseJSON['error'])
            $("#spinner").css("display", "none")
            $("#page-darken").css("opacity", "0")
            closeModal('add-existing-asset-modal');
            openModal('error-modal')
        }
    });
}

function deleteAsset(assetID) {
    parts = window.location.href.split('/')
    modelID = parts[parts.length - 1]
    
    data = {
        "model": modelID,
        "asset": assetID
    }

    $("#spinner").css("display", "block")
    $("#page-darken").css("opacity", "1")

    $.ajax({
        url: "/api/model/asset",
        type: "DELETE",
        contentType: 'application/json',
        data: JSON.stringify(data),
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
            assetID = response["id"]
            addAssetData = {
                "model": modelID,
                "asset": assetID
            }
            $.ajax({
                url: "/api/model/asset",
                type: "POST",
                contentType: 'application/json',
                data: JSON.stringify(addAssetData),
                success: function(response) {
                    $("#spinner").css("display", "none")
                    $("#page-darken").css("opacity", "0")
                    closeModal('add-git-asset-modal');
                    window.location.reload();
                },
                error: function(response) {
                    console.log(response)
                    $("#log-container").text(response.responseJSON['error'])
                    $("#spinner").css("display", "none")
                    $("#page-darken").css("opacity", "0")
                    closeModal('add-git-asset-modal');
                    openModal('error-modal')
                }
            });
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
            assetID = response["id"]
            addAssetData = {
                "model": modelID,
                "asset": assetID
            }
            $.ajax({
                url: "/api/model/asset",
                type: "POST",
                contentType: 'application/json',
                data: JSON.stringify(addAssetData),
                success: function(response) {
                    $("#spinner").css("display", "none")
                    $("#page-darken").css("opacity", "0")
                    closeModal('add-file-asset-modal');
                    window.location.reload();
                },
                error: function(response) {
                    console.log(response)
                    $("#log-container").text(response.responseJSON['error'])
                    $("#spinner").css("display", "none")
                    $("#page-darken").css("opacity", "0")
                    closeModal('add-file-asset-modal');
                    openModal('error-modal')
                }
            });
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

function dropdownSearchToggle() {
    document.getElementById("dropdown-dropdown").classList.toggle("show");
}

function filterFunction() {
    var input, filter, div, table, tbody, trs, tds, i;
    input = document.getElementById("dropdown-search");
    filter = input.value.toLowerCase();
    div = document.getElementById("dropdown-dropdown");
    table = document.getElementById("dropdown-table");
    tbody = table.getElementsByTagName("tbody")[0];
    trs = tbody.getElementsByTagName("tr");
    for (i = 1; i < trs.length; i++) {
        tds = trs[i].getElementsByTagName("td");
        objName = tds[0].innerText.toLowerCase();
        objID = tds[1].innerText.toLowerCase();
        if (objName.indexOf(filter) > -1 || objID.indexOf(filter) > -1) {
            trs[i].style.display = "";
        } else {
            trs[i].style.display = "none";
        }
    }
}

function downloadModel(url) {
    const a = document.createElement('a')
    a.href = url
    a.download = url.split('/').pop()
    document.body.appendChild(a)
    a.click()
    document.body.removeChild(a)
}

function syncGitAsset(assetID) {
    $("#spinner").css("display", "block")
    $("#page-darken").css("opacity", "1")

    $.ajax({
        url: "/api/asset/git/sync/" + assetID,
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
