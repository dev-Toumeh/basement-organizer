const NotificationTypeError = "error";
const NotificationTypeSuccess = "success";
const NotificationTypeWarning = "warning";
const NotificationTypeInfo = "";
/** @typedef {(typeof NotificationTypeError | typeof NotificationTypeSuccess | typeof NotificationTypeWarning | NotificationTypeInfo | undefined)} NotificationType
 * 
 * Possible notification types: "error", "success", "warning", "" (info), or undefined. */

/**
 * @typedef {Object} ServerNotification  
 * @property {string} message - The message text.
 * @property {NotificationTypeError | NotificationTypeSuccess | NotificationTypeWarning | NotificationTypeInfo | undefined} notificationType
 * @property {number | undefined} duration - How long is should display before removal. Default is 2000 (2 seconds).
 * @property {number | undefined} id - Add Id to HTML element: <div id="notification-{id}</div>. Will be random by default */

//document.getElementById('registerButton').addEventListener('click', function(event) {
//    var password = document.getElementById('password').value;
//    var confirmPassword = document.getElementById('password-confirm').value;
//    if (password !== confirmPassword) {
//        document.getElementById('response').innerText = 'Passwords do not match';
//        event.preventDefault(); // Prevent form from submitting
//    }
//});

//document.querySelectorAll("a")
//    .forEach(e => {
//        e.addEventListener('htmx:beforeRequest', stopRequestToSamePath)
//    })
//
//function stopRequestToSamePath(/** @type Event*/ event){
//    console.log(event.detail.pathInfo.requestPath, document.location.pathname)
//    if (event.detail.pathInfo.requestPath === document.location.pathname) {
//        event.preventDefault()
//    }
//}
//
function editItem(id) {
    const form = document.getElementById(`item-form-${id}`);
    const inputs = form.querySelectorAll('input:not([name="id"]):not([name="picture"])');
    inputs.forEach(input => input.removeAttribute('readonly'));

    document.getElementById(`image-preview-${id}`).style.display = 'none';
    document.getElementById(`image-input-${id}`).style.display = 'block';

    form.querySelector('button[onclick^="editItem"]').style.display = 'none';
    form.querySelector('button[type="submit"]').style.display = 'inline';
    form.querySelector('button[onclick^="cancelEdit"]').style.display = 'inline';
}

function cancelEdit(id) {
    const form = document.getElementById(`item-form-${id}`);
    const inputs = form.querySelectorAll('input:not([name="id"]):not([name="picture"])');
    inputs.forEach(input => {
        input.setAttribute('readonly', true);
        input.value = input.defaultValue;
    });

    document.getElementById(`image-preview-${id}`).style.display = 'block';
    document.getElementById(`image-input-${id}`).style.display = 'none';
    document.getElementById(`picture-${id}`).value = '';

    form.querySelector('button[onclick^="editItem"]').style.display = 'inline';
    form.querySelector('button[type="submit"]').style.display = 'none';
    form.querySelector('button[onclick^="cancelEdit"]').style.display = 'none';
}

/** Callback function that handles Error responses to display for the user.
 * Is registered as an eventListener
 * @param {CustomEvent<HtmxResponseInfo>} resInfo */
function errorResponseCallback(resInfo) {
    var notificationType
    var responseMessage
    var showNotification

    const status = resInfo.detail.xhr.status;

    // Show all error responses above 400 with a snackbar notification
    switch (true) {
        case status === 501:
            showNotification = true
            notificationType = NotificationTypeWarning;
            break;

        case status >= 400:
            showNotification = true
            notificationType = NotificationTypeError;
            console.log(resInfo);
            break;

        default:
            showNotification = false
            break;
    }

    if (showNotification) {
        resInfo.detail.isError = false;

        if (resInfo.detail.serverResponse !== "") {
            responseMessage = resInfo.detail.serverResponse;
        } else {
            responseMessage = resInfo.detail.xhr.statusText;
        }

        createAndShowNotification(responseMessage, notificationType);
    }

}

/** createAndShowNotification is for showing errors and information on updates.
 * @param {string} text - Message to display.
 * @param {NotificationTypeError | NotificationTypeInfo | NotificationTypeSuccess | NotificationTypeWarning | undefined } notificationType
 * @param {number | undefined} duration - How long is should display before removal. Default is 2000 (2 seconds).
 * @param {number | undefined} id - Add Id to HTML element: <div id="notification-{id}</div>. Will be random by default. */
function createAndShowNotification(text, notificationType, duration = 2000, id) {
    const snackbar = createNotification(text, notificationType, id);
    const snackbarElements = document.getElementById("notification-container");
    snackbarElements.appendChild(snackbar);

    showNotification(snackbar.id, duration)
}

/**
 * Creates a snackbar element but doesn't append it to the notification-container. 
 *
 * @param {string} text - The message to be displayed in the snackbar.
 * @param notificationType {NotificationTypeError | NotificationTypeInfo | NotificationTypeSuccess | NotificationTypeWarning | undefined } Can be "error" or undefined for info.
 * @param {string} [id] - Optional unique identifier for the snackbar. If not provided, a random ID is generated.
 * @returns {HTMLDivElement} The created snackbar element. */
function createNotification(text, notificationType, id) {
    const snackbar = document.createElement('div');
    const p = document.createElement("p");
    snackbar.appendChild(p)

    snackbar.className = "notification noshow " + notificationType;

    var snackbarId;
    if (id === undefined || id === "") {
        snackbarId = Math.round((Math.random() * 100000)).toString();
    } else {
        snackbarId = id;
    }

    snackbar.id = "notification-" + snackbarId;
    p.textContent = text;

    switch (notificationType) {
        case NotificationTypeError:
            console.error(text);
            break;
        case NotificationTypeInfo:
            console.info(text);
            break;
        case NotificationTypeWarning:
            console.warn(text);
            break;
        default:
            console.log(text);
            break;
    }
    return snackbar;
}

/** showNotification creates notification snackbar with id and automatically removes after duration.
 * @param id {string} HTML Element id.
 * @param text {string | undefined} Message to display.
 * @param duration {number | undefined} How long is should display before removal. Default is 2000 (2 seconds). */
function showNotification(id, duration = 2000) {
    let currentNotificationCount = document.querySelectorAll("div.notification").length;
    console.log("notification count: ", currentNotificationCount);

    // Add warning notification for too many notifications.
    if (currentNotificationCount > 10) {
        let warnNotificationId = 999999;
        let warnNotification = document.getElementById("notification-" + warnNotificationId.toString());

        if (!warnNotification) {
            warnNotification = createNotification("Over 10 notifications", NotificationTypeWarning, warnNotificationId);
            let notificationElements = document.getElementById("notification-container");
            notificationElements.prepend(warnNotification);

            setTimeout(() => {
                warnNotification.className = warnNotification.className.replace("noshow", "show");
            }, 50);

            setTimeout(function() {
                warnNotification.className = warnNotification.className.replace("show", "noshow");
                removeNotificationAfter(warnNotification.id, 210);
            }, duration);
        }
    }

    var notification = document.getElementById(id);
    if (notification === null) {
        console.error("no notification to show. ID: ", id);
    }

    console.log("notification added: ", notification)
    setTimeout(() => {
        notification.className = notification.className.replace("noshow", "show");
    }, 50);


    setTimeout(function() {
        notification.className = notification.className.replace("show", "noshow");
        removeNotificationAfter(id, 210);
    }, duration);
}

/** removeNotification removes snackbar HTML element.
 * @param id {string} HTML Element id. */
function removeNotification(id) {
    let bar = document.getElementById(id);
    if (bar !== null) {
        bar.remove();
    }
}

function removeAllNotifications() {
    document.getElementById("notification-container").innerHTML = "";
}

/** removeNotificationAfter removes snackbar HTML element after a certain duration.
 * @param id {string} HTML Element id. 
 * @param duration {number | undefined} How long to wait for removal. */
function removeNotificationAfter(id, duration) {
    setTimeout(() => removeNotification(id), duration);
}

/** noResponseCallback handles notification if requests don't respond.
 * @param resInfo {HmtxResponseInfo} */
function noResponseCallback(resInfo) {
    console.log(resInfo.requestConfig);
    createAndShowNotification("No response. Server down?", "error");
}

/** 
 * Callback function to show snackbar notifications triggered by the server.
 * Handles ServerNotificationEvents triggered with the "HX-Trigger" response header from the server.
 * <https://htmx.org/headers/hx-trigger/>
 *
 * @param {CustomEvent<{value: ServerNotification[]}>} evt - The event object containing details about the server notifications. */
function serverNotificationsCallback(evt) {
    serverNotifications(evt.detail.value);
}

/** Creates notifications from list.
 * @param {ServerNotification[]} notifications */
function serverNotifications(notifications) {
    for (let i = 0; i < notifications.length; i++) {
        createAndShowNotification(notifications[i].message, notifications[i].type, notifications[i].duration, notifications[i].id);
    }
}

/** Callback for header field "notification".
 * Will create notification with the value of that field.
 *
 * @param {CustomEvent<HtmxResponseInfo>} event */
function serverNotificationsFromHeaderCallback(event) {
    const hasNotificationEvent = event.detail.requestConfig.headers.hasOwnProperty("notification");
    if (hasNotificationEvent) {
        let serverNotification = event.detail.requestConfig.headers.notification;
        //debugger;
        let message = serverNotification.message ? serverNotification.message : "";
        let type = serverNotification.type ? serverNotification.type : "";
        let duration = serverNotification.duration ? parseInt(serverNotification.duration) : undefined;
        let id = serverNotification.id ? parseInt(serverNotification.id) : undefined;
        createAndShowNotification(message, type, duration, id);
    }
    if (event.detail.requestConfig.headers.hasOwnProperty("ServerNotificationEvents")) {
        serverNotifications(event.detail.requestConfig.headers.ServerNotificationEvents);
    }
}

/** Callback for handling checkboxes when they are changed. */
function checkboxChangedCallback(event) {
    if (event.target.classList.contains('move-checkbox') || event.target.classList.contains('delete-checkbox')) {
        // @TODO: Is this still used or needed?
        handleCheckboxChange();
    }
    if (event.target.matches("[type='checkbox']")) {
        const checkbox = event.target;
        if (checkbox.checked) {
            checkbox.setAttribute("checked", "");
        } else {
            checkbox.removeAttribute("checked");
        }
    }
}

var init = false;
function registerCallbackEventListener() {
    if (init) {
        return;
    }
    document.body.addEventListener('htmx:beforeSwap', errorResponseCallback);
    document.body.addEventListener('htmx:afterSwap', serverNotificationsFromHeaderCallback);
    document.body.addEventListener('htmx:sendError', noResponseCallback);
    document.body.addEventListener('ServerNotificationEvents', serverNotificationsCallback);
    document.body.addEventListener('change', checkboxChangedCallback);
    htmx.on('htmx:beforeHistorySave', removeAllNotifications);
    init = true;
}

console.log("registerCallbackEventListener executed")
//htmx.logAll();
htmx.onLoad(registerCallbackEventListener)


// responsible of disable checkBoxes if other type was checked
const MoveCheckboxes = document.getElementsByClassName('move-checkbox');
const DeleteCheckboxes = document.getElementsByClassName('delete-checkbox');

function handleCheckboxChange() {
    const moveChecked = document.querySelector('.move-checkbox:checked');
    const deleteChecked = document.querySelector('.delete-checkbox:checked');

    toggleCheckboxes(MoveCheckboxes, !!deleteChecked);
    toggleCheckboxes(DeleteCheckboxes, !!moveChecked);
}

function toggleCheckboxes(checkboxes, disable) {
    for (let i = 0; i < checkboxes.length; i++) {
        checkboxes[i].disabled = disable;
    }
}


let selectedIds = [];

function handleMoveItems(event) {
    const checkedCheckboxes = document.querySelectorAll('input.move-checkbox:checked');

    selectedIds = [];
    checkedCheckboxes.forEach((checkbox) => {
        selectedIds.push(checkbox.name);
    });

    console.log(selectedIds);
}

function triggerPaginationClickEvent(formID, page) {
    document.getElementById(formID + '-current-page').value = page;
    htmx.trigger("#" + formID, "paginationclick");
}

function toggleNav(el) {
  const nav = document.querySelector('.nav');
  nav.classList.toggle('nav-active');
  el.classList.toggle('hamburger--open');
}

