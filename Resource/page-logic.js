const ENV_PAGE_BTN = "env-page-btn";
const STRINGS_PAGE_BTN = "strings-page-btn";
const ASCENSIONS_PAGE_BTN = "ascensions-page-btn";
const TAOISTS_PAGE_BTN = "taoists-page-btn";

let buttons = [ENV_PAGE_BTN, STRINGS_PAGE_BTN, ASCENSIONS_PAGE_BTN, TAOISTS_PAGE_BTN];


let mainElement = document.getElementById("content-container");

for (let element of document.getElementsByClassName("page-switcher")) {
    element.addEventListener("click", onPageSwap);
}

function onPageSwap(event) {

    for (let buttonId of buttons) {
        document.getElementById(buttonId).classList.remove("selected");
    }

    event.target.classList.add("selected");
    const id = event.target.id;

    switch (id) {
        case ENV_PAGE_BTN:
            setEnvPage(id);
            break;
        case STRINGS_PAGE_BTN:
            setStringsPage(id);
            break;
        case ASCENSIONS_PAGE_BTN:
            setAscensionsPage(id);
            break;
        case TAOISTS_PAGE_BTN:
            setTaoistsPage(id);
            break;
    }
}

async function setPage(url, id) {
    const response = await fetch(url)

    if(response.status !== 200) {
        alert(`Failed to fetch ${id}:  ${response.statusText}`);
        return;
    }

    mainElement.innerHTML = await response.text();
}

function setEnvPage(id) { let _ = setPage("./env-page.html", id); }

function setStringsPage(id) { let _ = setPage("./strings-page.html", id); }

function setAscensionsPage(id) { let _ = setPage("./ascensions-page.html", id); }

function setTaoistsPage(id) { let _ = setPage("./taoists-page.html", id); }