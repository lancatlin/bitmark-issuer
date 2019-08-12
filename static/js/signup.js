let form = document.getElementById("signup_form")
form.onsubmit = function(ev) {
    signup(form)
    ev.preventDefault()
}
function signup(form) {
    let email = form.email.value
    let password = form.password.value
    let username = form.name.value
    let xhttp = new XMLHttpRequest()
    xhttp.open("post", "/signup", true)
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
                    error_msg.innerHTML = this.responseText
                    break
                case 500:
                    // server went wrong
                    alert(this.responseText)
                    break
            }
        }
    }
    xhttp.setRequestHeader("Content-Type", "application/x-www-form-urlencoded")
    xhttp.send(`email=${email}&password=${password}&name=${username}`)
}