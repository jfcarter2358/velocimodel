// import material.js
// import theme.js
// import modal.js
// import user_menu.js
// import status.js
// import data.js

var tagify

$(document).ready(
    function() {
        var input = document.getElementById('tag-input'),
        tagify = new Tagify(input, {
            id: 'release_tags',
        })
        var coll = document.getElementsByClassName("collapsible");
        var i;

        for (i = 0; i < coll.length; i++) {
            coll[i].addEventListener("click", function() {
                this.classList.toggle("active");
                var content = this.nextElementSibling;
                if (content.style.maxHeight){
                content.style.maxHeight = null;
                } else {
                content.style.maxHeight = content.scrollHeight + "px";
                } 
            });
        }
    }
)

function downloadRelease(url) {
    const a = document.createElement('a')
    a.href = url
    a.download = url.split('/').pop()
    document.body.appendChild(a)
    a.click()
    document.body.removeChild(a)
}
