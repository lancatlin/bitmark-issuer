let form = document.getElementById("login_form")
form.onsubmit = function(ev) {
    login(form)
    ev.preventDefault()
}
function login(form) {
    let email = form.email.value
    let password = form.password.value
    let xhttp = new XMLHttpRequest()
    xhttp.open("post", "/login", true)
    xhttp.onreadystatechange = function () {
        if (this.readyState == 4) {
            switch (this.status) {
                case 200:
                    // login success
                    window.location = "/"
                    break
                case 401:
                    // login fatal
                    let error_msg = document.getElementById("error_msg")
                    error_msg.hidden = false
                    break
                case 500:
                    // server went wrong
                    alert(this.responseText)
                    break
            }
        }
    }
    xhttp.setRequestHeader("Content-Type", "application/x-www-form-urlencoded")
    xhttp.send(`email=${email}&password=${password}`)
}