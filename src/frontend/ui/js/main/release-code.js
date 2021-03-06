// import material.js
// import theme.js
// import modal.js
// import user_menu.js
// import status.js
// import data.js

require.config({ paths: { vs: '../../../static/js/monaco-editor/min/vs' } });
var editor;
var model

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

            releaseJSONString = $("#release-json").text()
            releaseJSONFormatted = JSON.stringify(JSON.parse(releaseJSONString), null, 4);

            // Editor Config
            model = monaco.editor.createModel(releaseJSONFormatted, undefined, monaco.Uri.file('release.json'))

            editor = monaco.editor.create(document.getElementById('edit-container'), {
                model: model,
                minimap: {
                    enabled: false
                },
                automaticLayout: true,
                readOnly: true
            });
        }
    )
)
