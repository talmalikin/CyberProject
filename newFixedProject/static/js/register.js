document.getElementById('login-link').addEventListener('click', () => {
  window.location.href = "../index.html";
});

function validateUsername(username) {
  const re = /^[a-zA-Z0-9_]{3,}$/; // אותיות, מספרים, קו תחתון, מינימום 3 תווים
  return re.test(username);
}

function validatePassword(password) {
  return password.length >= 8; // מינימום 8 תווים
}

function showMessage(msg, isError = true) {
  let messageDiv = document.getElementById('message');
  if (!messageDiv) {
    messageDiv = document.createElement('div');
    messageDiv.id = 'message';
    document.body.insertBefore(messageDiv, document.body.firstChild);
  }
  messageDiv.style.color = isError ? 'red' : 'green';
  messageDiv.textContent = msg;
}

document.getElementById('register-btn').addEventListener('click', async () => {
  const username = document.getElementById('username').value.trim();
  const password = document.getElementById('password').value.trim();

  if (!username || !password) {
    showMessage('Please fill all fields');
    return;
  }
  if (!validateUsername(username)) {
    showMessage('Username must be at least 3 characters and contain only letters, numbers, or underscores');
    return;
  }
  if (!validatePassword(password)) {
    showMessage('Password must be at least 8 characters');
    return;
  }

  try {
    const res = await fetch('/api/register', {
      method: 'POST',
      headers: {'Content-Type': 'application/json'},
      body: JSON.stringify({username, password})
    });

    if (res.ok) {
      showMessage('Registration successful! Please login.', false);
      setTimeout(() => {
        window.location.href = "../index.html";
      }, 1500);
    } else {
      const errorText = await res.text();
      showMessage(errorText);
    }
  } catch (err) {
    showMessage('Network error');
  }
});
