// import material.js
// import theme.js
// import modal.js
// import user_menu.js
// import status.js
// import data.js

var releases

function getReleases() {
    $.ajax({
        url: "/api/release",
        type: "GET",
        success: function (result) {
            releases = result
        }
    });
}

function render() {
    var prefix = $("#search").val();
    var tempReleases = JSON.parse(JSON.stringify(releases));
    prefix = prefix.toLowerCase();
    if (prefix != "") {
        for (var idx = tempReleases.length - 1; idx >= 0; idx--) {
            if (tempReleases[idx].name.toLowerCase().indexOf(prefix) == -1) {
                tempReleases.splice(idx, 1)
            }
        }
    }
    var table = document.getElementById('releases-table');
    var tableHTMLString = '<tr><th class="table-title w3-medium velocimodel-text-blue"><span class="table-title-text">Release Name</span></th><th class="table-title w3-medium velocimodel-text-blue"><span class="table-title-text">Last Updated</span></th><th class="table-title w3-medium velocimodel-text-blue"><span class="table-title-text">Language</span></th><th class="table-title w3-medium velocimodel-text-blue"><span class="table-title-text">Version</span></th><th class="table-title w3-medium velocimodel-text-blue"><span class="table-title-text">Tags</span></th><th><div class="w3-round w3-button velocimodel-green">Create New</div></th></tr>' +
        tempReleases.map(function (release) {
            return '<tr>' +
                '<td>' + release.name + '</td>' +
                '<td>' + release.updated + '</td>' +
                '<td>' + release.language + '</td>' +
                '<td>' + release.version + '</td>' +
                '<td>' + release.tags.map(function (tag) {
                    return '<span class="w3-tag w3-round velocimodel-green release-tag">' + tag + '</span>'
                }).join('') +
                '</td>' +
                '<td class="table-link-cell">' +
                '<a href="/ui/release/' + release.id + '" class="table-link-link w3-right-align light theme-text" style="float:right;margin-right:16px;"><i class="fa-solid fa-link"></i></a>' +
                '</td>' +
                '</tr>'
        }).join('');

    table.innerHTML = tableHTMLString;
}

$(document).ready(
    function() {
        getReleases()
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

function createReleaseFromSnapshot(snapshotID) {
    $("#spinner").css("display", "block")
    $("#page-darken").css("opacity", "1")

    $.ajax({
        url: "/api/snapshot/" + snapshotID + "/release",
        type: "POST",
        success: function(response) {
            $("#spinner").css("display", "none")
            $("#page-darken").css("opacity", "0")
            releaseID = response["id"]
            closeModal('create-from-snapshot-modal');
            window.location.assign('/ui/release/' + releaseID);
        },
        error: function(response) {
            console.log(response)
            $("#log-container").text(response.responseJSON['error'])
            $("#spinner").css("display", "none")
            $("#page-darken").css("opacity", "0")
            closeModal('create-from-snapshot-modal');
            openModal('error-modal')
        }
    });
}
