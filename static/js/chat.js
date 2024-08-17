const host = 'http://192.168.1.14:8080/';
const url = host + `get_messages?chatID=${chatID}`;
console.log(url);
fetch(url)
    .then(response => {
        if (!response.ok) {
            throw new Error('Network response was not ok');
        }
        return response.json();
    })
    .then(data => {
        console.log(data);
        if (data.messages != null) {
            data.messages.forEach(mess => {
                console.log(mess.id);
                addMessage(mess.user, mess.time, mess.message);
            });
        }
    })
    .catch(error => {
        alert(error);
    });

document.cookie = `user_login=${encodeURIComponent(sessionStorage.getItem('user_login'))}; path=/`;
const conn = new WebSocket('ws://192.168.1.14:8080/ws?chatID=' + chatID);
conn.onmessage = function (event) {
    var message = JSON.parse(event.data);
    console.log(message.id);
    addMessage(message.user, message.time, message.message);
};

function addMessage(user, time, text) {
    var chatbox = document.getElementById("chatbox");
    chatbox.innerHTML += `<div class="message">
            <span class="user">${decodeURIComponent(user)}</span>
            <span class="time">${time}</span>
            <div class="text">${text}</div>
            <div class="actions">
                <button class="edit-btn">Изменить</button>
                <button class="delete-btn">Удалить</button>
            </div>
        </div>`;
    chatbox.scrollTop = chatbox.scrollHeight; // Scroll to the bottom
}

function sendMessage() {
    var input = document.getElementById("message");
    if (input.value.length === 0) {
        return;
    }
    const now = new Date();
    const day = String(now.getDate()).padStart(2, '0');
    const month = String(now.getMonth() + 1).padStart(2, '0');
    const year = now.getFullYear();
    const hours = String(now.getHours()).padStart(2, '0');
    const minutes = String(now.getMinutes()).padStart(2, '0');
    const formattedDate = `${day}.${month}.${year} ${hours}:${minutes}`;
    console.log(`in sendmessage: ${formattedDate}`);

    const messageWithTime = {
        time: formattedDate,
        message: input.value
    };

    conn.send(JSON.stringify(messageWithTime));
    input.value = "";
}

// Add event listener for Enter key
document.getElementById("message").addEventListener("keydown", function (event) {
    if (event.key === "Enter") {
        event.preventDefault();
        sendMessage();
    }
});
