const host = 'http://192.168.1.14:8080/';
console.log("hello");

function toggleForm(formType) {
    if (formType === 'login') {
        document.getElementById('login-form').style.display = 'block';
        document.getElementById('register-form').style.display = 'none';
    } else if (formType === 'register') {
        document.getElementById('login-form').style.display = 'none';
        document.getElementById('register-form').style.display = 'block';
    }
}

document.getElementById('login_button').addEventListener('click', function () {
    if (document.getElementById('login-nickname').value.length === 0) {
        alert("введите ник");
        return;
    }
    if (document.getElementById('login-password').value.length === 0) {
        alert("введите пароль");
        return;
    }
    fetch(host + 'login', {
        method: 'POST',
        headers: {
            'Content-Type': 'application/json'
        },
        body: JSON.stringify({
            login: document.getElementById('login-nickname').value,
            password: document.getElementById('login-password').value
        })
    }).then(response => {
        if (response.ok) {
            sessionStorage.setItem('user_login', decodeURIComponent(document.getElementById('login-nickname').value));
            window.location.href = host + 'welcome';
        } else {
            return response.text().then(text => {
                throw new Error(text);
            });
        }
    }).catch(error => {
        alert(error.message);
    });
});


document.getElementById('register_button').addEventListener('click', function () {
    if (document.getElementById('register-nickname').value.length === 0) {
        alert("введите ник");
        return;
    }
    if (document.getElementById('register-password').value.length === 0) {
        alert("введите пароль");
        return;
    }
    fetch(host + 'register', {
        method: 'POST',
        headers: {
            'Content-Type': 'application/json'
        },
        body: JSON.stringify({
            login: document.getElementById('register-nickname').value,
            password: document.getElementById('register-password').value
        })
    }).then(response => {
        if (response.ok) {
            alert("вы успешно зарегестрировались!");
            window.location.href = host;
        } else {
            return response.text().then(text => {
                throw new Error(text);
            });
        }
    }).catch(error => {
        alert(error.message);
    });
});


