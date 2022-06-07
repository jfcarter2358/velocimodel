
// import material.js
// import theme.js
// import modal.js
// import user_menu.js
// import status.js
// import data.js

function saveUser() {
    parts = window.location.href.split('/')
    userID = parts[parts.length - 1]

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
        url: "/api/user/" + userID,
        type: "PUT",
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
