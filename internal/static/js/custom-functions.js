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
            createAndShowSnackbar(resInfo.detail.serverResponse, "error");
        } else {
            createAndShowSnackbar(resInfo.detail.xhr.statusText, "error");
        }
        //resInfo.detail.target = htmx.find("#snackbars");
        //resInfo.detail.target.setAttribute("hx-swap", "afterend");
    }
}

const SnackbarTypeError = "error";
const SnackbarTypeInfo = "";

/** createAndShowSnackbar is for showing errors and information on updates.
 * @param text {string} Message to display.
 * @param type {SnackbarTypeError | SnackbarTypeInfo | undefined } Can be "error" or undefined for info.
 * @param duration {number | undefined} How long is should display before removal. Default is 2000 (2 seconds). */
function createAndShowSnackbar(text, type, duration = 2000) {
    const newElement = document.createElement('div');

    newElement.className = 'snackbar noshow ' + type;
    newElement.id = "snackbar-" + Math.round((Math.random() * 100)).toString();

    const snackbarsElement = document.getElementById('snackbars');

    snackbarsElement.appendChild(newElement);

    showSnackbarWithId(newElement.id, text, duration)
}


/** showSnackbarWithId creates notification snackbar with id and automatically removes after duration.
 * @param id {string} HTML Element id.
 * @param text {string | undefined} Message to display.
 * @param duration {number | undefined} How long is should display before removal. Default is 2000 (2 seconds). */
function showSnackbarWithId(id, text, duration = 2000) {
    var snackbar = document.getElementById(id);

    setTimeout(() => {
        snackbar.className = snackbar.className.replace("noshow", "show");
        snackbar.textContent = text;
    }, 50);


    setTimeout(function() {
        snackbar.className = snackbar.className.replace("show", "noshow");
        removeSnackbarAfter(id, 210);
    }, duration);
}

/** removeSnackbar removes snackbar HTML element.
 * @param snackbarId {string} HTML Element id. */
function removeSnackbar(snackbarId) {
    bar = document.getElementById(snackbarId);
    bar.remove();
}

/** removeSnackbar removes snackbar HTML element after a certain duration.
 * @param snackbarId {string} HTML Element id. 
 * @param duration {number | undefined} How long to wait for removal. */
function removeSnackbarAfter(snackbarId, duration) {
    setTimeout(() => removeSnackbar(snackbarId), duration);
}

/** noResponseCallback handles notification if requests don't respond.
 * @param resInfo {HmtxResponseInfo} */
function noResponseCallback(resInfo) {
    console.log(resInfo.requestConfig);
    createAndShowSnackbar("No response. Server down?", "error");
}

function registerCallbackEventListener() {
    document.body.addEventListener('htmx:beforeSwap', errorResponseCallback);
    document.body.addEventListener('htmx:sendError', noResponseCallback);
}

console.log("registerCallbackEventListener executed")
htmx.onLoad(registerCallbackEventListener)
