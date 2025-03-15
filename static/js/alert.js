'use strict';

let alertPlaceholder = document.getElementById('alertPlaceholder');
let alertCount = 0;

function appendAlert(message, type = "success") {
    if (Object.is(alertPlaceholder, null)) {
        alertPlaceholder = document.getElementById('alertPlaceholder');
    }

    const wrapper = document.createElement('div');
    wrapper.id = "alert-" + alertCount++;
    wrapper.className = `alert alert-${type} alert-dismissible`;
    wrapper.role = "alert";

    const messageDiv = document.createElement('div');
    messageDiv.textContent = message;
    wrapper.append(messageDiv);

    const dismissButton = document.createElement('button');
    dismissButton.type = "button";
    dismissButton.className = "btn-close";
    dismissButton.setAttribute("data-bs-dismiss", "alert");
    dismissButton.setAttribute("aria-label", "Close");
    dismissButton.onclick = function () {
        document.getElementById(wrapper.id).remove();
    };
    wrapper.append(dismissButton);

    alertPlaceholder.append(wrapper);
}
