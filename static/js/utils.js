'use strict';

let alertPlaceholder = document.getElementById('alertPlaceholder');
let alertCount = 0;

function appendAlert(message, type = "success") {
    if (Object.is(alertPlaceholder, null)) {
        alertPlaceholder = document.getElementById('alertPlaceholder');
    }

    const count = alertCount++;

    const wrapper = document.createElement('div');
    wrapper.id = "alert-" + count;
    wrapper.className = `alert alert-${type} alert-dismissible`;
    wrapper.role = "alert";
    wrapper.setAttribute("aria-live", type === "success" ? "polite" : "assertive");

    const messageDiv = document.createElement('div');
    messageDiv.id = `alert-message-${count}`;
    messageDiv.textContent = message;
    wrapper.setAttribute("aria-describedby", messageDiv.id);
    wrapper.append(messageDiv);

    const dismissButton = document.createElement('button');
    dismissButton.type = "button";
    dismissButton.className = "btn-close";
    dismissButton.setAttribute("data-bs-dismiss", "alert");
    dismissButton.setAttribute("aria-label", "Close alert");
    dismissButton.onclick = function () {
        let currentAlert = document.getElementById(wrapper.id);
        if (currentAlert) {
            currentAlert.remove();
        }
    };
    wrapper.append(dismissButton);

    alertPlaceholder.append(wrapper);
}

let modal = null;

function showAddHostModal() {
    if (!modal) {
        modal = new bootstrap.Modal(document.getElementById('addHostModal'));
    }
    modal.show();
}

function formatAndValidateMAC(input) {
    // Remove all non-hexadecimal characters
    let value = input.value.replace(/[^a-fA-F0-9]/g, '').toUpperCase();

    // Format the MAC address with colons
    let formattedValue = value.match(/.{1,2}/g)?.join(':') || '';

    // Limit to 17 characters (6 pairs of hex digits + 5 colons)
    if (formattedValue.length > 17) {
        formattedValue = formattedValue.slice(0, 17);
    }

    input.value = formattedValue;

    // Pattern is already enforced by the formatting above.
    // Only check left is the length.
    if (formattedValue.length === 17) {
        input.setCustomValidity("");
    } else {
        input.setCustomValidity("MAC Address must be 17 characters long");
    }
}

function validateHostname(input) {
    const hostname = input.value;

    if (/[^a-zA-Z0-9.-]/.test(hostname)) {
        input.setCustomValidity("Hostname can only contain letters, numbers, dots, and hyphens");
        return;
    }
    if (hostname.length > 253) {
        input.setCustomValidity("Hostname must be less than 253 characters");
        return;
    }
    if (hostname.length < 1) {
        input.setCustomValidity("Hostname cannot be empty");
        return;
    }
    if (hostname.startsWith('.') || hostname.endsWith('.')) {
        input.setCustomValidity("Hostname cannot start or end with a dot");
        return;
    }

    const labels = hostname.split('.');
    for (const label of labels) {
        if (label.length > 63) {
            input.setCustomValidity("Each label in the hostname must be less than 63 characters");
            return;
        }
        if (label.length < 1) {
            input.setCustomValidity("Each label in the hostname must be at least 1 character");
            return;
        }
        if (label.startsWith('-') || label.endsWith('-')) {
            input.setCustomValidity("Labels cannot start or end with a hyphen");
            return;
        }
    }

    input.setCustomValidity("");
}

function updateHostStatus(statuses) {
    for (let status of statuses) {
        const addressElement = document.getElementById(status.mac + ".Address");
        if (!addressElement) {
            console.warn(`No address element found for host: ${status.mac}`);
            continue;
        }
        addressElement.innerHTML = status.online ? "🟢 " + status.address : "🔴 " + status.address;
        if (status.error) {
            appendAlert(`Failed to fetch status for host ${status.mac}: ${status.error}`, "warning");
        }
    }
}

hostStatus();
// Update host status every 30 seconds
setInterval(hostStatus, 30000);
