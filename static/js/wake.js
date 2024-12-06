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
