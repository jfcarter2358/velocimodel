// import material.js
// import theme.js
// import modal.js
// import user_menu.js
// import status.js
// import data.js

require.config({ paths: { vs: '../../../static/js/monaco-editor/min/vs' } });
var editor;
var model

function saveAsset() {
    parts = window.location.href.split('/')
    basePath = ''
    basePath = ''
if (parts[3] != 'ui') {
        basePath = "/" + parts[3]
    }
    assetID = parts[parts.length - 2]

    data = JSON.parse(editor.getValue())

    $("#spinner").css("display", "block")
    $("#page-darken").css("opacity", "1")

    $.ajax({
        url: basePath + "/api/asset/" + assetID,
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

$(document).ready(
    require(['vs/editor/editor.main'],
        function () {
            monaco.editor.setTheme('vs');
            monaco.editor.defineTheme('velocimodelLightTheme', {
                base: 'vs',
                inherit: true,
                rules: [{ background: 'ffffff' }],
                colors: {
                    'editor.foreground': '#000000',
                    'editor.background': '#ffffff',
                    'editorCursor.foreground': '#f4f6f5',
                    'editor.lineHighlightBackground': '#00000020',
                    'editorLineNumber.foreground': '#838995',
                    'editor.selectionBackground': '#83899530',
                    'editor.inactiveSelectionBackground': '#00000015',
                    'editorWidget.background': '#ffffff',
                    'editorWidget.border': '#e1e7e4',
                    'editor.selectionBackground': '#83899530',
                    'editor.inactiveSelectionBackground': '#838995015',
                    'editorHoverWidget.background': '#ffffff',
                    'editorHoverWidget.border': 'e1e7e4'
                }
            });
            monaco.editor.defineTheme('velocimodelDarkTheme', {
                base: 'vs-dark',
                inherit: true,
                rules: [{ background: '333333' }],
                colors: {
                    'editor.foreground': '#000000',
                    'editor.background': '#333333',
                    'editorCursor.foreground': '#f4f6f5',
                    'editor.lineHighlightBackground': '#00000020',
                    'editorLineNumber.foreground': '#838995',
                    'editor.selectionBackground': '#83899530',
                    'editor.inactiveSelectionBackground': '#00000015',
                    'editorWidget.background': '#333333',
                    'editorWidget.border': '#4b4b4b',
                    'editor.selectionBackground': '#83899530',
                    'editor.inactiveSelectionBackground': '#838995015',
                    'editorHoverWidget.background': '#333333',
                    'editorHoverWidget.border': '4b4b4b'
                }
            });
            if ($(".main").hasClass("light")) {
                monaco.editor.setTheme('velocimodelLightTheme');
            } else {
                monaco.editor.setTheme('velocimodelDarkTheme');
            }

            assetJSONString = $("#asset-json").text()
            assetJSONFormatted = JSON.stringify(JSON.parse(assetJSONString), null, 4);

            // Editor Config
            model = monaco.editor.createModel(assetJSONFormatted, undefined, monaco.Uri.file('asset.json'))

            editor = monaco.editor.create(document.getElementById('edit-container'), {
                model: model,
                minimap: {
                    enabled: false
                },
                automaticLayout: true
            });
        }
    )
)
