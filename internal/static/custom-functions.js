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

