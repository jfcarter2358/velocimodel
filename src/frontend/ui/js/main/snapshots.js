// import material.js
// import theme.js
// import modal.js

var snapshots

function getSnapshots() {
    $.ajax({
        url: "/script/api/snapshot",
        type: "GET",
        success: function (result) {
            snapshots = result
        }
    });
}

function render() {
    var prefix = $("#search").val();
    var tempSnapshots = JSON.parse(JSON.stringify(snapshots));
    prefix = prefix.toLowerCase();
    if (prefix != "") {
        for (var idx = tempSnapshots.length - 1; idx >= 0; idx--) {
            if (tempSnapshots[idx].name.toLowerCase().indexOf(prefix) == -1) {
                tempSnapshots.splice(idx, 1)
            }
        }
    }
    var table = document.getElementById('snapshots-table');
    var tableHTMLString = '<tr><th class="table-title w3-medium velocimodel-text-blue"><span class="table-title-text">Snapshot Name</span></th><th class="table-title w3-medium velocimodel-text-blue"><span class="table-title-text">Last Updated</span></th><th class="table-title w3-medium velocimodel-text-blue"><span class="table-title-text">Language</span></th><th class="table-title w3-medium velocimodel-text-blue"><span class="table-title-text">Version</span></th><th class="table-title w3-medium velocimodel-text-blue"><span class="table-title-text">Tags</span></th><th><div class="w3-round w3-button velocimodel-green">Create New</div></th></tr>' +
        tempSnapshots.map(function (snapshot) {
            return '<tr>' +
                '<td>' + snapshot.name + '</td>' +
                '<td>' + snapshot.updated + '</td>' +
                '<td>' + snapshot.language + '</td>' +
                '<td>' + snapshot.version + '</td>' +
                '<td>' + snapshot.tags.map(function (tag) {
                    return '<span class="w3-tag w3-round velocimodel-green snapshot-tag">' + tag + '</span>'
                }).join('') +
                '</td>' +
                '<td class="table-link-cell">' +
                '<a href="/ui/snapshot/' + snapshot.id + '" class="table-link-link w3-right-align dark theme-text" style="float:right;margin-right:16px;"><i class="fa-solid fa-link"></i></a>' +
                '</td>' +
                '</tr>'
        }).join('');

    table.innerHTML = tableHTMLString;
}

$(document).ready(
    function() {
        getSnapshots()
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

function createSnapshotFromModel(modelID) {
    $("#spinner").css("display", "block")
    $("#page-darken").css("opacity", "1")

    $.ajax({
        url: "/script/api/model/" + modelID + "/snapshot",
        type: "POST",
        success: function(response) {
            $("#spinner").css("display", "none")
            $("#page-darken").css("opacity", "0")
            snapshotID = response["id"]
            closeModal('create-from-model-modal');
            window.location.assign('/ui/snapshot/' + snapshotID);
        },
        error: function(response) {
            console.log(response)
            $("#log-container").text(response.responseJSON['error'])
            $("#spinner").css("display", "none")
            $("#page-darken").css("opacity", "0")
            closeModal('create-from-model-modal');
            openModal('error-modal')
        }
    });
}
