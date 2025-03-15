async function wake(macAddr) {
    try {
        const response = await fetch("/api/" + macAddr);

        const responseBody = await response.json();

        if (response.ok) {
            alert("Send magic packet");
        } else {
            alert("Error: " + responseBody.reason);
        }
    } catch (error) {
        console.error(error.message);
    }
}

async function wakeCustom() {
    const inputCustomMAC = document.getElementById("custom-mac-input");
    wake(inputCustomMAC.value);
}
