
// import material.js
// import theme.js
// import modal.js
// import user_menu.js
// import status.js
// import data.js

var users

function getUsers() {
    parts = window.location.href.split('/')
    if (parts[0] != 'ui') {
        basePath = "/" + parts[3]
    }
    $.ajax({
        url: basePath + "/api/user",
        type: "GET",
        success: function (result) {
            users = result
        }
    });
}

function render() {
    var prefix = $("#search").val();
    var tempUsers = JSON.parse(JSON.stringify(users));
    prefix = prefix.toLowerCase();
    if (prefix != "") {
        for (var idx = tempUsers.length - 1; idx >= 0; idx--) {
            if (tempUsers[idx].name.toLowerCase().indexOf(prefix) == -1) {
                tempUsers.splice(idx, 1)
            }
        }
    }
    var table = document.getElementById('users-table');
    var tableHTMLString = '<tr><th class="table-title w3-medium velocimodel-text-blue"><span class="table-title-text">User Name</span></th><th class="table-title w3-medium velocimodel-text-blue"><span class="table-title-text">Last Updated</span></th><th class="table-title w3-medium velocimodel-text-blue"><span class="table-title-text">Type</span></th><th class="table-title w3-medium velocimodel-text-blue"><span class="table-title-text">Tags</span></th><th><div class="w3-round w3-button velocimodel-green">Create New</div></th></tr>' +
        tempUsers.map(function (user) {
            return '<tr>' +
                '<td>' + user.username + '</td>' +
                '<td>' + user.roles + '</td>' +
                '<td>' + user.groups + '</td>' +
                '<td>' + user.email + '</td>' +
                '<td class="table-link-cell">' +
                '<a href="{{ .base_path }}/ui/user/' + user.id + '" class="table-link-link w3-right-align light theme-text" style="float:right;margin-right:16px;"><i class="fa-solid fa-link"></i></a>' +
                '</td>' +
                '</tr>'
        }).join('');

    table.innerHTML = tableHTMLString;
}

function addUser() {
    parts = window.location.href.split('/')
    if (parts[0] != 'ui') {
        basePath = "/" + parts[3]
    }

    data = {
        "username": $("#users-add-username").val(),
        "password": $("#users-add-password").val(),
        "given_name": $("#users-add-given-name").val(),
        "family_name": $("#users-add-family-name").val(),
        "id": "",
        "roles": $("#users-add-roles").val(),
        "groups": $("#users-add-groups").val(),
        "email": $("#users-add-email").val(),
        "reset_token": "",
        "reset_token_create_date": "",
        "created": "",
        "updated": ""
    }

    $("#spinner").css("display", "block")
    $("#page-darken").css("opacity", "1")

    $.ajax({
        url: basePath + "/api/user",
        type: "POST",
        contentType: 'application/json',
        data: JSON.stringify(data),
        success: function(response) {
            $("#spinner").css("display", "none")
            $("#page-darken").css("opacity", "0")
            userID = response['id']
            closeModal('users-add-modal');
            window.location.assign(basePath + '/ui/user/' + userID);
        },
        error: function(response) {
            closeModal('users-delete-modal');
            console.log(response)
            $("#log-container").text(response.responseJSON['error'])
            $("#spinner").css("display", "none")
            $("#page-darken").css("opacity", "0")
            openModal('error-modal')
        }
    });
}

function openDeleteModal(username, id) {
    $("#user-delete-username").text(username)
    $("#user-delete-id").text(id)
    openModal("users-delete-modal");
}

function deleteUser() {
    parts = window.location.href.split('/')
    if (parts[0] != 'ui') {
        basePath = "/" + parts[3]
    }

    id = $("#user-delete-id").text()
    data = [id]

    $("#spinner").css("display", "block")
    $("#page-darken").css("opacity", "1")

    $.ajax({
        url: basePath + "/api/user",
        type: "DELETE",
        contentType: 'application/json',
        data: JSON.stringify(data),
        success: function(response) {
            $("#spinner").css("display", "none")
            $("#page-darken").css("opacity", "0")
            closeModal('users-delete-modal');
            window.location.reload()
        },
        error: function(response) {
            closeModal('users-delete-modal');
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
        getUsers()
    }
)
