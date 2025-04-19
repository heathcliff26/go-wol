'use strict';

async function wake(macAddr, name = "") {
    const displayName = name != "" ? name + " (" + macAddr + ")" : macAddr

    try {
        const response = await fetch("/api/v1/wake/" + macAddr);

        const responseBody = await response.json();

        if (response.ok) {
            appendAlert("Send magic packet to " + displayName);
        } else {
            appendAlert("Failed to send magic packet to " + displayName + " : " + responseBody.reason, "warning");
        }
    } catch (error) {
        console.error(error.message);
        appendAlert("Failed to send magic packet to " + displayName, "danger");
    }
}

async function wakeCustom() {
    const inputCustomMAC = document.getElementById("custom-mac-input");
    wake(inputCustomMAC.value);
}

async function addHost() {
    const name = document.getElementById('hostName').value;
    const macAddr = document.getElementById('macAddress').value;

    modal.hide();
    try {
        const response = await fetch('/api/v1/hosts/'+macAddr+"/"+name, {
            method: 'PUT',
        });

        const responseBody = await response.json();

        if (response.ok) {
            appendAlert(`Added host ${name}`);
            location.reload();
        } else {
            appendAlert(`Failed to add host: ${responseBody.reason}`, "warning");
        }
    } catch (error) {
        console.error(error.message);
        appendAlert("Failed to add host " + name, "danger");
    }
}

async function deleteHost(macAddr, name) {
    if (!confirm(`Are you sure you want to delete ${name}?`)) {
        return;
    }

    try {
        const response = await fetch(`/api/v1/hosts/${macAddr}`, {
            method: 'DELETE',
        });

        const responseBody = await response.json();

        if (response.ok) {
            appendAlert(`Deleted host ${name}`);
            location.reload();
        } else {
            appendAlert(`Failed to delete host: ${responseBody.reason}`, "warning");
        }
    } catch (error) {
        console.error(error.message);
        appendAlert("Failed to delete host " + name, "danger");
    }
}
