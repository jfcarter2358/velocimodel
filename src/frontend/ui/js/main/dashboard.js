// import material.js
// import data.js
// import theme.js

function LoadModels() {
    $.ajax({
        type: "GET",
        url: '/api/model',
        success: function(response) {
            for(var i = 0; i < response.length; i++) {
                console.log(response[i])
                $("#models-table").append(
                    '<tr>' +
                    '<td>' +
                    '<a href="/ui/model/' +
                    response[i]['name'] +
                    '" class="w3-button dark theme-hover-light w3-border-bottom theme-border-light">' +
                    response[i]['name'] +
                    '</a>' +
                    '</td>' +
                    '<td>' +
                    response[i]['updated'] +
                    '</td>' +
                    '<td>' +
                    response[i]['language'] +
                    '</td>' +
                    '</tr>'
                )
            }
        },
        error: function(response) {
            console.log(response)
        }
    });
}

// $(document).ready(
//     function() {
//         LoadModels()
//     }
// )
