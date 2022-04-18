var tagify

$(document).ready(
    function() {
        var input = document.getElementById('tag-input'),
        tagify = new Tagify(input, {
            id: 'test1',
        })
    }
)

function getTags() {
    console.log(document.getElementById('tag-input').value)
}
