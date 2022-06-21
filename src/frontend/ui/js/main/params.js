
// import material.js
// import theme.js
// import modal.js
// import user_menu.js
// import status.js
// import data.js

function saveParams() {
    parts = window.location.href.split('/')
    if (parts[0] != 'ui') {
        basePath = "/" + parts[3]
    }

    data = {}
    $("#params-card").children("input").each(function() {
        key = $(this).attr('id')
        key = key.slice("params-".length)
        data[key] = this.value
    })

    $("#spinner").css("display", "block")
    $("#page-darken").css("opacity", "1")

    $.ajax({
        url: basePath + "/api/param",
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
