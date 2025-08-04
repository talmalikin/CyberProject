document.getElementById('go-register').addEventListener('click', () => {
  window.location.href = "html/register.html";
});

document.getElementById('login-btn').addEventListener('click', async () => {
  const username = document.getElementById('username').value.trim();
  const password = document.getElementById('password').value.trim();
  const msgEl = document.getElementById('message');

  function showMessage(text) {
    msgEl.textContent = text;
  }

  if (!username) {
    showMessage('Please enter your username.');
    return;
  }
  if (!username.match(/^[a-zA-Z0-9_]{3,30}$/)) {
    showMessage('Username must be 3-30 characters, letters, numbers or underscores only.');
    return;
  }
  if (!password) {
    showMessage('Please enter your password.');
    return;
  }
  if (password.length < 8) {
    showMessage('Password must be at least 8 characters long.');
    return;
  }

  try {
    const res = await fetch('/api/login', {
      method: 'POST',
      headers: {'Content-Type': 'application/json'},
      body: JSON.stringify({username, password})
    });

    if (res.ok) {
      showMessage('');
      window.location.href = "html/vault.html";
    } else {
      showMessage('Login failed. Please check your credentials.');
    }
  } catch (error) {
    showMessage('Network error. Please try again later.');
  }
});
