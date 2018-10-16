function bindNewBookButtion() {
    let newBookButton = document.querySelector("div[for='control'] .new-book")
    let uploaderSection = document.querySelector("div[for=upload-books]")
    newBookButton.addEventListener("click", function () {
        let display = uploaderSection.style.display;
        let nextDisplay = display === "none" ? "block" : "none";
        uploaderSection.style.display = nextDisplay;
    })
}

function bindNewFolderForm() {
    let form = document.querySelector("div[for='control'] form")
    let newListButton = document.querySelector("[for=control] button.new-list")
    let nameInput = form.querySelector("input")
    newListButton.addEventListener("click", function () {
        form.style.display = "block"
        nameInput.focus()
    })
    form.addEventListener("submit", function (evt) {


        evt.preventDefault()
        data = new FormData(this)
        name = data.get("name")
        name = name.trim()

        let request = new XMLHttpRequest();
        request.open('POST', newListEndPoint, true);
        request.setRequestHeader('Content-Type', 'application/json; charset=UTF-8');
        request.send(JSON.stringify({
            name
        }));

        request.onload = function () {
            if (request.status >= 200 && request.status < 400) {
                data = JSON.parse(request.responseText)
                form.reset();
                put2Lists(data)
                form.style.display = "none"
            } else {
                console.error("server error!")
            }
        }

    })
}

function put2Lists(data) {
    let {
        name,
        id
    } = data;
    let divEle = document.createElement("div")
    divEle.className = "list-child"
    divEle.innerHTML = `<a href="/u/explorer/${id}"><i class="material-icons">folder</i></a><span>${name}<span>`

    let container = document.querySelector("[for=show-lists] [role=lists]").appendChild(divEle)
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