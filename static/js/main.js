const host = 'http://192.168.1.14:8080/';
const curr_login = decodeURIComponent(sessionStorage.getItem('user_login'));

function generateChatID(user1, user2) {
    // Сортируем имена по алфавиту, чтобы порядок не влиял на chatID
    return [user1, user2].sort().join('-');
}

const url = host + `get_users?user_login=${sessionStorage.getItem('user_login')}`;
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
                window.location.href = `${host}chat?chatID=${chatID}`;
            };

            li.appendChild(button);
            userList.appendChild(li);
        });
    })
    .catch(error => {
        alert(error);
    });
