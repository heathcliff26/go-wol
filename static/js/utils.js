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
