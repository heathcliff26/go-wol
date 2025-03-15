'use strict';

async function wake(macAddr, name = "") {
    const displayName = name != "" ? name + " (" + macAddr + ")" : macAddr

    try {
        const response = await fetch("/api/" + macAddr);

        const responseBody = await response.json();

        if (response.ok) {
            alert("Send magic packet to " + displayName);
        } else {
            alert("Failed to send magic packet to " + displayName + " : " + responseBody.reason);
        }
    } catch (error) {
        console.error(error.message);
    }
}

async function wakeCustom() {
    const inputCustomMAC = document.getElementById("custom-mac-input");
    wake(inputCustomMAC.value);
}
