
const stringsObserver = new MutationObserver((mutationsList) => {
    for (const mutation of mutationsList) {
        if (mutation.type === "childList") {
            console.log("New HTML content added:", mutation.addedNodes);

            for (let node of mutation.addedNodes) {
                if(node.nodeType !== 1) return;

                if (node.classList.contains("table-container")) {
                    const table = node.childNodes[1]

                    if("strings-table" === table.id) {
                        //populate table
                        populateStringTable();
                    }
                }
            }
        }
    }
});

stringsObserver.observe(document.body, {
    childList: true,
    subtree: true
});

function addRow(key, value) {
    const table = document.getElementById("strings-table-body");


    const tr = document.createElement("tr");

    tr.innerHTML = `
        <td><input type="text" class="key-input" placeholder="Key"></td>
        <td><input type="text" class="value-input" placeholder="Value"></td>
    `;


    const keyInput = tr.querySelector(".key-input");
    const valueInput = tr.querySelector(".value-input");

    keyInput.value = key;
    valueInput.value = value;

    table.appendChild(tr);
}

async function saveStrings() {

    const table = document.getElementById("strings-table-body");

    let jsonEntries = [];

    for (const row of table.children) {
        const key = row.querySelector(".key-input")?.value.trim();
        const value = row.querySelector(".value-input")?.value.trim();

        if (key) // skipping empty values
            jsonEntries.push({key, value });
    }

    const response = await fetch('/save-strings', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify(jsonEntries),
    });

    if (403 === response.status) {
        Toasty.showToast('Relog', 'address is not authenticated.', '', 2000);
        setTimeout(function(){
            window.location = "./login.html";
        }, 2000)
        return
    }

    if(!response.ok) {
        Toasty.showToast('Relog', 'Something went wrong, if the error persists contact your developer.', '', 2000);
        return
    }

    Toasty.showToast('Success', "Modified translations successfully.", '');
}

async function populateStringTable() {
    const response = await fetch("./strings");

    if(403 === response.status){
        Toasty.showToast('Relog', 'address is not authenticated.', '', 2000);
        setTimeout(function(){
            window.location = "./login.html";
        }, 2000)
        return
    }

    if(!response.ok) {
        Toasty.showToast('Relog', 'Something went wrong, if the error persists contact your developer.', '', 2000);
        return
    }

    let stringInfo = Object.entries(await response.json()).map(([key, value]) => ({ key, value }));

    for (let stringkvp of stringInfo) {
        addRow(stringkvp.key, stringkvp.value.value);
    }
}

document.addEventListener("click", function (event) {

    const element = event.target;

    switch (element.id) {
        case "add-string-button":
            addRow("", "");
            break;
        case "save-string-button":
            saveStrings();
            break;
    }
});