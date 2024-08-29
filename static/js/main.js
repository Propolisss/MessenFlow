let host = 'http://';

fetch('/get_socket')
    .then(response => {
        if (!response.ok) {
            throw new Error('Network response was not ok');
        }
        return response.text();
    })
    .then(data => {
        console.log('data:', data);
        host += data;
        console.log('host:', host);
        run();
    })
    .catch(error => {
        alert(error);
    });

function run() {

    const curr_login = decodeURIComponent(sessionStorage.getItem('user_login'));

    function generateChatID(user1, user2) {
        return [user1, user2].sort().join('-');
    }

    const url = host + `get_users?user_login=${sessionStorage.getItem('user_login')}`;
    console.log('url:', url);
    fetch(url)
        .then(response => {
            if (!response.ok) {
                throw new Error('Ошибка: ' + response.status);
            }
            return response.json();
        })
        .then(data => {
            const userList = document.getElementById('user-list');
            data.users.forEach(user => {
                const li = document.createElement('li');
                li.textContent = user;

                const button = document.createElement('button');
                button.textContent = 'Написать сообщение';
                button.onclick = () => {
                    const chatID = generateChatID(curr_login, user);
                    sessionStorage.setItem('friend_login', decodeURIComponent(user));
                    window.location.href = `${host}chat?chatID=${chatID}`;
                };

                li.appendChild(button);
                userList.appendChild(li);
            });
        })
        .catch(error => {
            alert(error);
        });
}
