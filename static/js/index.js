document.getElementById('go-register').addEventListener('click', () => {
  window.location.href = "html/register.html";
});

document.getElementById('login-btn').addEventListener('click', async () => {
  const username = document.getElementById('username').value.trim();
  const password = document.getElementById('password').value.trim();

  if (!username || !password) {
    alert('Please fill all fields');
    return;
  }

  const res = await fetch('/api/login', {
    method: 'POST',
    headers: {'Content-Type': 'application/json'},
    body: JSON.stringify({username, password})
  });

  if (res.ok) {
    alert('Login successful');
    window.location.href = "html/vault.html";
  } else {
    alert(await res.text());
  }
});
