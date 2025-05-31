const observer = new MutationObserver((mutationsList) => {
    for (const mutation of mutationsList) {
        if (mutation.type === "childList") {
            console.log("New HTML content added:", mutation.addedNodes);

            for (let node of mutation.addedNodes) {
                if(node.nodeType !== 1) return;

                if (node.classList.contains("table-container")) {
                    const table = node.childNodes[1]

                    if("env-table" === table.id) {
                        //populate table
                        populateTable();
                    }
                }
            }
        }
    }
});

observer.observe(document.body, {
    childList: true,
    subtree: true
});

async function populateTable() {

    const response = await fetch("env");

    if(response.status === 403){
        Toasty.showToast('Relog', 'address is not authenticated.', '', 2000);
        setTimeout(function(){
            window.location = "./login.html";
        }, 2000)
        return;
    }

    if(!response.ok) {
        Toasty.showToast('Relog', 'Something went wrong, if the error persists contact your developer.', '', 2000);
        return
    }

    const envRows = await response.json();

    for(let env  of envRows)
        addEnvRow(env.type, env.key, env.value);
}


function addEnvRow(type = "string", key = "", value = "") {
    const tbody = document.getElementById("env-table-body");

    const tr = document.createElement("tr");

    tr.innerHTML = `
        <td>
            <select class="type-select">
                <option value="string">string</option>
                <option value="int">int</option>
                <option value="bool">bool</option>
                <option value="float">float</option>
            </select>
        </td>
        <td><input type="text" class="key-input" placeholder="Key"></td>
        <td><input type="text" class="value-input" placeholder="Value"></td>
        <td>
            <div class="row-actions">
                <button class="action-btn delete-btn row-delete" title="Delete">🗑️</button>
            </div>
        </td>
    `;

    // Set initial values
    const select = tr.querySelector(".type-select");
    const keyInput = tr.querySelector(".key-input");
    const valueInput = tr.querySelector(".value-input");

    if (["string", "int", "bool", "float"].includes(type)) {
        select.value = type;
    }

    keyInput.value = key;
    valueInput.value = value;

    tbody.appendChild(tr);
}

async function saveEnv() {
    const table = document.getElementById("env-table-body");

    let jsonEntries = [];

    for (const row of table.children) {
        const type = row.querySelector(".type-select")?.value;
        const key = row.querySelector(".key-input")?.value.trim();
        const value = row.querySelector(".value-input")?.value.trim();

        if (key) // skipping empty values
            jsonEntries.push({ type, key, value });
    }

    const response = await fetch('/save-env', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify(jsonEntries),
    });

    if (response.status === 403) {
        Toasty.showToast('Relog', 'address is not authenticated.', '', 2000);
        return;
    }

    if(!response.ok) {
        Toasty.showToast('Relog', 'Something went wrong, if the error persists contact your developer.', '', 2000);
        return;
    }

    Toasty.showToast('Success', 'Saved env values successfully.', '', 2000);
}

function removeRow(deleteBtnElement) {
    const row = deleteBtnElement.closest('tr'); // or '.row' or 'tr'
    if (row) row.remove();
}

document.addEventListener("", (event) => {
    console.log(event.target);
})

document.addEventListener("click", function (event) {

    const element = event.target;

    switch(element.id) {
        case "add-env-button":
            addEnvRow();
            break;
        case "save-env-button":
            saveEnv();
            break;
    }

    if(element.classList.contains("row-delete")){
        removeRow(element);
    }
})