document.getElementById('login-link').addEventListener('click', () => {
  window.location.href = "../index.html";
});

document.getElementById('register-btn').addEventListener('click', async () => {
  const username = document.getElementById('username').value.trim();
  const password = document.getElementById('password').value.trim();

  if (!username || !password) {
    alert('Please fill all fields');
    return;
  }

  try {
    const res = await fetch('/api/register', {
      method: 'POST',
      headers: {'Content-Type': 'application/json'},
      body: JSON.stringify({username, password})
    });

    if (res.ok) {
      alert('Registration successful! Please login.');
      window.location.href = "../index.html";
    } else {
      const errorText = await res.text();
      alert(errorText);
    }
  } catch (err) {
    alert('Network error');
  }
});
