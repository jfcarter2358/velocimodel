
// import material.js
// import theme.js
// import modal.js
// import user_menu.js
// import status.js
// import data.js

function saveSecrets() {
    
    data = {}
    $("#secrets-card").children("input").each(function() {
        key = $(this).attr('id')
        key = key.slice("secrets-".length)
        data[key] = this.value
    })

    console.log(data)

    $("#spinner").css("display", "block")
    $("#page-darken").css("opacity", "1")

    $.ajax({
        url: "/api/secret",
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
