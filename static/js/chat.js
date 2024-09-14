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
        socket = data;
        console.log('host:', host);
        run();
    })
    .catch(error => {
        alert(error);
    });

function run() {
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
                    console.log(mess.id, mess.type);
                    addMessage(mess.id, mess.user, mess.time, mess.message);
                });
            }
        })
        .catch(error => {
            alert(error);
        });

    document.cookie = `user_login=${encodeURIComponent(sessionStorage.getItem('user_login'))}; path=/`;

    const conn = new WebSocket('ws://' + socket + 'ws?chatID=' + chatID);
    conn.onmessage = function (event) {
        var message = JSON.parse(event.data);
        console.log(message.id, message.type);
        if (message.type === 'add') {
            addMessage(message.id, message.user, message.time, message.message);
        } else if (message.type === 'update') {
            const messageElement = document.querySelector(`.message[data-message-id="${message.id}"]`);
            const textElement = messageElement.querySelector('.text');
            textElement.innerText = message.message;
        } else if (message.type === 'delete') {
            const messageElement = document.querySelector(`.message[data-message-id="${message.id}"]`);
            if (messageElement) {
                messageElement.remove();
            }
        } else if (message.type === 'status') {
            console.log(message.user, message.online);
            updateStatus(message.online);
        }
    };


    function addMessage(id, user, time, text) {
        var chatbox = document.getElementById("chatbox");
        var actionsHtml = '';
        var messageClass = '';

        if (user === sessionStorage.getItem('user_login')) {
            actionsHtml = `
            <div class="actions">
                <button class="edit-btn">Изменить</button>
                <button class="delete-btn">Удалить</button>
            </div>`;
            messageClass = 'own-message';
        }

        chatbox.innerHTML += `
        <div class="message ${messageClass}" data-message-id="${id}">
            <span class="user">${decodeURIComponent(user)}</span>
            <span class="time">${time}</span>
            <div class="text">${text}</div>
            ${actionsHtml}
        </div>`;

        chatbox.scrollTop = chatbox.scrollHeight;
    }

    let currentlyEditingForm = null;
    let currentHtmlInner = null;
    let currentMessageId = null;

    function updateMessage(messageId, currentMessage) {
        console.log(`in updateMessage: ${messageId},${currentMessage}`);

        fetch(host + `update_message?message_id=${messageId}&new_text=${currentMessage}`, {
            method: 'PUT',
        })
            .then(response => {
                if (!response.ok) {
                    throw new Error(`database error: can not update message with messageId: ${messageId}`);
                }
                const updateMessage = {
                    type: 'update',
                    id: parseInt(messageId),
                    message: currentMessage,
                };
                conn.send(JSON.stringify(updateMessage));
            })
            .catch(error => {
                alert(error);
            });
    }

    document.getElementById('sendMessage').addEventListener('click', sendMessage);


// TODO: fix the logic when a new message is added
    document.getElementById('chatbox').addEventListener('click', function (event) {
        if (event.target.className === 'delete-btn') {
            const messageElement = event.target.closest('.message');
            const messageId = messageElement.getAttribute('data-message-id');
            console.log(`in deleteButtonClick: ${messageId}`);
            deleteMessage(messageId);
            messageElement.remove();
        } else if (event.target.className === 'edit-btn') {
            console.log(`in edit-btn: ${currentlyEditingForm}`);
            if (currentlyEditingForm) {
                const newText = currentlyEditingForm.querySelector('input').value;
                console.log(newText, currentlyEditingForm);
                currentlyEditingForm.remove();
                currentHtmlInner.innerHTML = '';
                currentHtmlInner.innerText = newText;
                updateMessage(currentMessageId, newText);
            }

            const messageElement = event.target.closest('.message');
            const messageId = messageElement.getAttribute('data-message-id');
            const textElement = messageElement.querySelector('.text');
            const currentText = textElement.innerText;


            const editForm = document.createElement('form');
            editForm.innerHTML = `
            <input type="text" value="${currentText}" />
            <button type="submit">Save</button>
        `;

            textElement.innerHTML = '';
            textElement.appendChild(editForm);

            currentlyEditingForm = editForm;
            currentHtmlInner = textElement;
            currentMessageId = messageId;

            editForm.addEventListener('submit', function (e) {
                e.preventDefault();
                const newText = editForm.querySelector('input').value;
                console.log(`new text: ${newText}`);
                currentlyEditingForm = null;
                currentHtmlInner = null;
                currentMessageId = null;
                updateMessage(messageId, newText);
            });
        }
    });

    function deleteMessage(messageId) {
        console.log(messageId);
        fetch(host + `delete_message?id=${messageId}`, {
            method: 'DELETE',
        })
            .then(response => {
                if (!response.ok) {
                    throw new Error(`database error: can not delete message with messageId: ${messageId}`);
                }
                const deleteMessage = {
                    type: 'delete',
                    id: parseInt(messageId),
                };
                conn.send(JSON.stringify(deleteMessage));
            })
            .catch(error => {
                alert(error);
            });
    }

    function sendMessage() {
        const input = document.getElementById("message");
        if (input.value.length === 0) {
            return;
        }
        if (currentlyEditingForm) {
            const newText = currentlyEditingForm.querySelector('input').value;
            console.log(newText, currentlyEditingForm);
            currentlyEditingForm.remove();
            currentHtmlInner.innerHTML = '';
            currentHtmlInner.innerText = newText;
            updateMessage(currentMessageId, newText);
            currentlyEditingForm = null;
            currentHtmlInner = null;
            currentMessageId = null;
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
            message: input.value,
            type: "add",
        };

        conn.send(JSON.stringify(messageWithTime));
        input.value = "";
    }

    document.getElementById("message").addEventListener("keydown", function (event) {
        if (event.key === "Enter") {
            event.preventDefault();
            sendMessage();
        }
    });
}

function updateStatus(online) {
    const chatTitle = document.getElementById('chat-title');
    const statusIndicator = document.createElement('span');
    statusIndicator.id = 'status-indicator';
    if (online) {
        statusIndicator.innerText = ' В сети';
        statusIndicator.classList.add('online');
        statusIndicator.classList.remove('offline');
    } else {
        statusIndicator.innerText = ' Не в сети';
        statusIndicator.classList.add('offline');
        statusIndicator.classList.remove('online');
    }
    //Удаляем существующий индикатор статуса, если он есть
    const existingStatusIndicator = chatTitle.querySelector('#status-indicator');
    if (existingStatusIndicator) {
        chatTitle.removeChild(existingStatusIndicator);
    }
    // Добавляем новый индикатор статуса
    chatTitle.appendChild(statusIndicator);
}