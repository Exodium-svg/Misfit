const ascensionsObserver = new MutationObserver((mutationsList) => {
    for (const mutation of mutationsList) {
        if (mutation.type === "childList") {
            console.log("New HTML content added:", mutation.addedNodes);

            for (let node of mutation.addedNodes) {
                if(node.nodeType !== 1) return;

                if (node.classList.contains("table-container")) {
                    const table = node.childNodes[1]

                    if("ascensions-table" === table.id) {
                        //populate table
                        populateAscensionsTable();
                    }
                }
            }
        }
    }
});

ascensionsObserver.observe(document.body, {
    childList: true,
    subtree: true
});


document.addEventListener("click", function (event) {

    const element = event.target;

    switch (element.id) {
        case "add-ascension-button":
            addAscensionRow(0, "", "");
            break;
        case "save-ascension-button":
            saveAscensions();
            break;
    }
});

let roles = [];


async function setRoles() {
    const response = await fetch("./get-roles");

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

    roles = await response.json();

    console.log(roles);
}

async function populateAscensionsTable() {
    await setRoles();

    const response =  await fetch("./get-ascensions");

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

    const ascensions = await response.json();

    console.log(ascensions);

    for (const ascension of ascensions) {
        addAscensionRow(ascension.required_level, ascension.role_id);
    }

}

function addAscensionRow(requiredLevel, roleId) {
    let role = {
        name:"",
        id:"",
    };
    const table = document.getElementById("ascensions-table");

    for (let i = 0; i < roles.length; i++) {
        if (roles[i] === roleId) {
            role = roles[i];
        }
    }

    const tr = document.createElement("tr");

    let options = ``

    for (let i = 0; i < roles.length; i++) {
        let localRole = roles[i]

        if(localRole.name === "@everyone") {
            continue;
        }

        if(localRole.id !== roleId)
            options += `<option value="${localRole.id}">${localRole.name}</option>`
        else
            options += `<option value="${localRole.id}" selected>${localRole.name}</option>`
    }

    tr.innerHTML = `
        <td><input type="number" class="key-input" placeholder="Key"></td>
        <td>
            <select class="type-select">
                ${options}
            </select>
        </td>
    `;


    const levelInput = tr.querySelector(".key-input");

    levelInput.value = requiredLevel;

    table.appendChild(tr);
}

async function saveAscensions() {
    const table = document.getElementById("ascensions-table");

    let jsonEntries = [];

    for (const row of table.children) {
        const levelInput = row.querySelector(".key-input");
        const select = row.querySelector(".type-select");

        if (!levelInput || !select) continue;

        const requiredLevel = parseInt(levelInput.value.trim(), 10);
        const selectedOption = select.options[select.selectedIndex];

        const roleId = selectedOption.value;
        const title = selectedOption.textContent;

        if (!isNaN(requiredLevel) && roleId) {
            jsonEntries.push({
                title,
                requiredLevel,
                roleId
            });
        }
    }

    const response = await fetch('/save-ascensions', {
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

    Toasty.showToast('Success', 'Saved ascensions successfully.', '', 2000);
}