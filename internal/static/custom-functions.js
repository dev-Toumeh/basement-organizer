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

//document.body.addEventListener('htmx:oobBeforeSwap', function(evt) {
//   console.log(evt) 
//});

document.body.addEventListener('htmx:beforeSwap', function(evt) {
    if(evt.detail.xhr.status === 404){
        // alert the user when a 404 occurs (maybe use a nicer mechanism than alert())
        //alert("Error: Could Not Find Resource");
    } else if(evt.detail.xhr.status === 400){
        //console.log(evt.detail);
        //alert("Bad request");
        //document.querySelector("#responses").innerHTML = evt.detail.xhr.status;
    } else if(evt.detail.xhr.status === 422){
        // allow 422 responses to swap as we are using this as a signal that
        // a form was submitted with bad data and want to rerender with the
        // errors
        //
        // set isError to false to avoid error logging in console
        evt.detail.shouldSwap = true;
        evt.detail.isError = false;
    } else if(evt.detail.xhr.status === 418){
        // if the response code 418 (I'm a teapot) is returned, retarget the
        // content of the response to the element with the id `teapot`
        evt.detail.shouldSwap = true;
        evt.detail.target = htmx.find("#teapot");
    }
    if(evt.detail.xhr.status >= 400){
        //showSnackbar(evt.detail.xhr.statusText + " " + evt.detail.serverResponse, 500000);
        console.log(evt);

        // Will still swap contents to show error notification
        evt.detail.shouldSwap = true;

        //evt.detail.target = htmx.find("#snackbars");
        //evt.detail.target.setAttribute("hx-swap", "afterend");
    }
});


function removeSnackbar(snackbarId) {
    bar = document.getElementById(snackbarId);
    bar.remove();
}

function removeSnackbarAfter(snackbarId, duration) {
    setTimeout(() => removeSnackbar(snackbarId), duration);
}

function showSnackbarWithId(id, text, duration) {
    if (duration === undefined) {
        duration = 2000;
    }

    var snackbar = document.getElementById(id);

    setTimeout(() => {
        snackbar.className = snackbar.className.replace("noshow", "show"); 
        snackbar.innerHTML = text;
    }, 50);


    setTimeout(function(){
        snackbar.className = snackbar.className.replace("show", "noshow"); 
       removeSnackbarAfter(id, 210);
    }, duration);

}

function showSnackbar(text, duration) {
    if (duration === undefined) {
        duration = 2000;
    }

    var x = document.getElementById("snackbar");

    if (currentSnackbarTimerId !== undefined) {
        clearTimeout(currentSnackbarTimerId);
        x.className = x.className.replace("show", ""); 
    }

    setTimeout(() => {
        x.className = "show";
        x.innerHTML = text;
    }, 50);


    currentSnackbarTimerId = setTimeout(function(){
        x.className = x.className.replace("show", ""); 
        currentSnackbarTimerId = undefined;
    }, duration);
}
