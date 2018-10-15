function bindNewBookButtion() {
    let newBookButton = document.querySelector("div[for='control'] .new-book")
    let uploaderSection = document.querySelector("div[for=upload-books]")
    newBookButton.addEventListener("click", function () {
        let display = uploaderSection.style.display;
        console.log(display);
        let nextDisplay = display === "none" ? "block" : "none";
        uploaderSection.style.display = nextDisplay;
    })
}

function bindNewFolderForm() {
    let form = document.querySelector("div[for='control'] form")

    form.addEventListener("submit", function (evt) {
        evt.preventDefault()
        data = new FormData(this)
        name = data.get("name")
        name = name.trim()

        let request = new XMLHttpRequest();
        request.open('POST', newListEndPoint, true);
        request.setRequestHeader('Content-Type', 'application/json; charset=UTF-8');
        request.send(JSON.stringify({name}));

        request.onload = function () {
            if (request.status >= 200 && request.status < 400) {
                // data = JSON.parse(request.responseText)
                put2Lists(data)

            } else {
                console.error("server error!")
            }
        }

    })
}

function put2Lists() {

}

function ready(fn) {
    if (document.attachEvent ? document.readyState === "complete" : document.readyState !== "loading") {
        fn();
    } else {
        document.addEventListener('DOMContentLoaded', fn);
    }
}

ready(function () {
    bindNewBookButtion()
    bindNewFolderForm()
})