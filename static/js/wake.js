'use strict';

async function wake(macAddr, name = "") {
    const displayName = name != "" ? name + " (" + macAddr + ")" : macAddr

    try {
        const response = await fetch("/api/" + macAddr);

        const responseBody = await response.json();

        if (response.ok) {
            appendAlert("Send magic packet to " + displayName);
        } else {
            appendAlert("Failed to send magic packet to " + displayName + " : " + responseBody.reason, "warning");
        }
    } catch (error) {
        console.error(error.message);
    }
}

async function wakeCustom() {
    const inputCustomMAC = document.getElementById("custom-mac-input");
    wake(inputCustomMAC.value);
}
