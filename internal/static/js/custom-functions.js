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
 * @param {HtmxResponseInfo}resInfo */
function errorResponseCallback(resInfo) {
    // Show all error responses above 400 with a snackbar notification
    if (resInfo.detail.xhr.status >= 400) {
        console.log(resInfo);
        resInfo.detail.isError = false;

        if (resInfo.detail.serverResponse !== "") {
            createAndShowNotification(resInfo.detail.serverResponse, "error");
        } else {
            createAndShowNotification(resInfo.detail.xhr.statusText, "error");
        }
        //resInfo.detail.target = htmx.find("#notification-container");
        //resInfo.detail.target.setAttribute("hx-swap", "afterend");
    }
}

const NotificationTypeError = "error";
const NotificationTypeSuccess = "success";
const NotificationTypeWarning = "warning";
const NotificationTypeInfo = "";

/** createAndShowNotification is for showing errors and information on updates.
 * @param text {string} Message to display.
 * @param notificationType {NotificationTypeError | NotificationTypeInfo | NotificationTypeSuccess | NotificationTypeWarning | undefined } Can be "error" or undefined for info.
 * @param duration {number | undefined} How long is should display before removal. Default is 2000 (2 seconds).
 * @param id {number | undefined} Add Id to HTML element: <div id="notification-{id}</div>. Will be random by default */
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
 * @returns {HTMLDivElement} The created snackbar element.
 */
function createNotification(text, notificationType, id) {
    const snackbar = document.createElement('div');

    snackbar.className = "notification noshow " + notificationType;

    var snackbarId;
    if (id === undefined || id === "") {
        snackbarId = Math.round((Math.random() * 100000)).toString(); 
    } else {
        snackbarId = id;
    }

    snackbar.id = "notification-" + snackbarId;
    snackbar.textContent = text;
    return snackbar;
}

/** showNotification creates notification snackbar with id and automatically removes after duration.
 * @param id {string} HTML Element id.
 * @param text {string | undefined} Message to display.
 * @param duration {number | undefined} How long is should display before removal. Default is 2000 (2 seconds). */
function showNotification(id, duration = 2000) {
    let currentNotificationCount = document.querySelectorAll("div.notification").length;
    console.log(currentNotificationCount);

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
 * @param {CustomEvent} evt - The event object containing details about the server notifications.
 * @param {Object} evt.detail - The event detail payload.
 * @param {Array} evt.detail.value - An array of notification objects.
 * @param {string} evt.detail.value[].message - The message to be displayed in the snackbar.
 * @param {string} evt.detail.value[].type - The type of notification (e.g., "success", "error").
 * @param {number} [evt.detail.value[].duration] - Optional duration for which the snackbar is displayed.
 * @param {string} [evt.detail.value[].id] - Optional unique identifier for the snackbar.
 */
function serverNotificationsCallback(evt){
    for (let i = 0; i < evt.detail.value.length; i++) {
        createAndShowNotification(evt.detail.value[i].message, evt.detail.value[i].type, evt.detail.value[i].duration, evt.detail.value[i].id);
    }
}

/** Callback for header field "notification".
 * Will create notification with the value of that field.
 * */
function serverNotificationsFromHeaderCallback(event) {
    const serverNotification = event.detail.requestConfig.headers.notification;
    if (serverNotification !== undefined) {
        //debugger;
        let message = serverNotification.message ? serverNotification.message : "";
        let type = serverNotification.type ? serverNotification.type : "";
        let duration = serverNotification.duration ? parseInt(serverNotification.duration) : undefined;
        let id = serverNotification.id ? parseInt(serverNotification.id) : undefined;
        createAndShowNotification(message, type, duration, id);
    }
}

function registerCallbackEventListener() {
    document.body.addEventListener('htmx:beforeSwap', errorResponseCallback);
    document.body.addEventListener('htmx:afterSwap', serverNotificationsFromHeaderCallback);
    document.body.addEventListener('htmx:sendError', noResponseCallback);
    document.body.addEventListener('ServerNotificationEvents', serverNotificationsCallback);
}

console.log("registerCallbackEventListener executed")
//htmx.logAll();
htmx.onLoad(registerCallbackEventListener)
