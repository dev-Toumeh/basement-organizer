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

// Callback function that handles Error responses
// Is registered as an eventListener
function errorResponseCallback(evt) {
    // Show all error responses above 400 with a snackbar notification
    if (evt.detail.xhr.status >= 400) {
        console.log(evt);
        evt.detail.isError = false;

        if (evt.detail.serverResponse !== "") {
            createAndShowSnackbar(evt.detail.serverResponse, "error");
        } else {
            createAndShowSnackbar(evt.detail.xhr.statusText, "error");
        }
        //evt.detail.target = htmx.find("#snackbars");
        //evt.detail.target.setAttribute("hx-swap", "afterend");
    }
}

function createAndShowSnackbar(text, type, duration) {
    const newElement = document.createElement('div');

    newElement.className = 'snackbar noshow ' + type;
    newElement.id = "snackbar-" + Math.round( (Math.random() * 100)).toString();

    const snackbarsElement = document.getElementById('snackbars');

    snackbarsElement.appendChild(newElement);

    showSnackbarWithId(newElement.id, text, duration)
    console.log('snack')
}

function showSnackbarWithId(id, text, duration) {
    if (duration === undefined) {
        duration = 2000;
    }

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

function removeSnackbar(snackbarId) {
    bar = document.getElementById(snackbarId);
    bar.remove();
}

function removeSnackbarAfter(snackbarId, duration) {
    setTimeout(() => removeSnackbar(snackbarId), duration);
}


function registerCallbackEventListener() {
    document.body.addEventListener('htmx:beforeSwap', errorResponseCallback );
}

console.log("registerCallbackEventListener executed")
htmx.onLoad(registerCallbackEventListener)
