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
