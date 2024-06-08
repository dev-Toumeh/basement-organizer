document.getElementById('registerButton').addEventListener('click', function(event) {
    var password = document.getElementById('password').value;
    var confirmPassword = document.getElementById('password-confirm').value;
    if (password !== confirmPassword) {
        document.getElementById('response').innerText = 'Passwords do not match';
        event.preventDefault(); // Prevent form from submitting
    }
});
