function bindNewBookButton() {
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
        let display = form.style.display;
        let nextDisplay = display === "none" ? "block" : "none";
        form.style.display = nextDisplay;

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
    divEle.classList.add("list-child", "selectable")
    divEle.dataset.id = id
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
    bindNewBookButton()
    bindNewFolderForm()
    bindSelection()
    bindDeleteButton()
})

function bindSelection() {
    function getEleRoot(ele) {
        let oldE = ele

        while (!ele.classList.contains('selectable')) {
            ele = ele.parentNode;
            if (ele.classList === undefined) {
                return oldE
            }
        }
        return ele
    }
    const selection = Selection.create({

        // Class for the selection-area
        class: 'selection',

        // All elements in this container can be selected
        containers: ['[for=show-lists', '[for=show-books]'],

        // The container is also the boundary in this case
        boundaries: ['[for=show-lists', '[for=show-books]'],

        onSelect(evt) {
            eleRoot = getEleRoot(evt.target)
            // Check if clicked element is already selected
            const selected = eleRoot.classList.contains('selected');

            // Remove class if the user don't pressed the control key or ⌘ key and the
            // current target is already selected
            if (!evt.originalEvent.ctrlKey && !evt.originalEvent.metaKey) {

                // Remove class from every element which is selected
                evt.selectedElements.forEach((s) => {
                    root = getEleRoot(s)
                    root.classList.remove('selected')
                });

                // Clear previous selection
                this.clearSelection();
            }

            if (!selected) {

                // Select element
                eleRoot.classList.add('selected');
                this.keepSelection();
            } else {

                // Unselect element
                eleRoot.classList.remove('selected');
                this.removeFromSelection(evt.target);
            }

            toggleMenu()

        },

        onStart(evt) {
            // Get elements which has been selected so far
            const selectedElements = evt.selectedElements;

            // Remove class if the user don't pressed the control key or ⌘ key
            if (!evt.originalEvent.ctrlKey && !evt.originalEvent.metaKey) {

                // Unselect all elements
                selectedElements.forEach((s) => {
                    root = getEleRoot(s)
                    root.classList.remove('selected')
                });

                // Clear previous selection
                this.clearSelection();
            }

        },

        onMove(evt) {

            // Get the currently selected elements and those
            // which where removed since the last selection.
            const selectedElements = evt.selectedElements;
            const removedElements = evt.changedElements.removed;

            // Add a custom class to the elements which where selected.
            selectedElements.forEach((s) => {
                root = getEleRoot(s)
                root.classList.add('selected')
            });

            // Remove the class from elements which where removed
            // since the last selection
            removedElements.forEach((s) => {
                root = getEleRoot(s)
                root.classList.remove('selected')
            });

            toggleMenu()

        },

        onStop() {
            this.keepSelection();
        }
    });
}

let button = document.querySelector("[for=control] button.remove")

function toggleMenu() {
    let selected = document.querySelectorAll(".selectable.selected")

    if (selected.length > 0) {
        button.style.display = "block";
    } else {
        button.style.display = "none";
    }
}

function bindDeleteButton() {
    button.addEventListener("click", function (evt) {
        evt.preventDefault()

        let lists = [].slice.call(document.querySelectorAll(".list-child.selectable.selected")).map(e => e.dataset.id)
        let books = [].slice.call(document.querySelectorAll(".book.selectable.selected")).map(e => e.dataset.id)


        let request = new XMLHttpRequest();
        request.open('DELETE', `/u/explorer/${currentList}`, true);
        request.setRequestHeader('Content-Type', 'application/json; charset=UTF-8');
        request.send(JSON.stringify({
            lists,
            books,
        }));

        request.onload = function () {
            location.reload()
        }
    })
}