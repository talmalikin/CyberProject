document.getElementById('toggle-add-btn').addEventListener('click', toggleAddFields);
document.querySelector('.hide-btn').addEventListener('click', toggleAddFields);
document.getElementById('add-entry-btn').addEventListener('click', addEntry);
document.querySelector('.logout-btn').addEventListener('click', logout);

function toggleAddFields() {
  const fields = document.getElementById('add-fields');
  const toggleBtn = document.getElementById('toggle-add-btn');
  if (fields.style.display === 'none' || fields.style.display === '') {
    fields.style.display = 'flex';
    toggleBtn.style.display = 'none';
  } else {
    fields.style.display = 'none';
    toggleBtn.style.display = 'inline-block';
    clearAddFields();
  }
}

function clearAddFields() {
  document.getElementById('website').value = '';
  document.getElementById('email').value = '';
  document.getElementById('password').value = '';
}

async function fetchEntries() {
  const res = await fetch('/api/vault');
  if (res.status === 401) {
    alert('Unauthorized. Please login.');
    window.location.href = "../index.html";
    return;
  }
  const data = await res.json();
  const container = document.getElementById('entries');
  container.innerHTML = '';
  data.forEach((item, index) => {
    const initials = item.website.slice(0, 2).toUpperCase();
    container.innerHTML += `
      <div class="entry">
        <div class="entry-left">
          <div class="icon-box">${initials}</div>
          <div class="entry-info">
            <div class="entry-line">
              <span><strong>${item.website}</strong></span>
            </div>
            <div class="entry-line">
              <span class="icon">👤</span><span>${item.email}</span>
            </div>
            <div class="entry-line password-box">
              <span class="icon">🔒</span>
              <input class="password-input" type="password" value="${item.password}" readonly />
              <button class="eye-btn" onclick="toggleVisibility(this)">👁️</button>
            </div>
          </div>
        </div>
        <span class="delete-btn" onclick="deleteEntry(${index})">🗑️</span>
      </div>`;
  });
}

async function addEntry() {
  const website = document.getElementById('website').value.trim();
  const email = document.getElementById('email').value.trim();
  const password = document.getElementById('password').value.trim();

  if (!website || !email || !password) {
    alert('Please fill all fields.');
    return;
  }

  const res = await fetch('/api/add', {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify({ website, email, password })
  });

  if (res.ok) {
    clearAddFields();
    toggleAddFields();
    fetchEntries();
  } else if (res.status === 401) {
    alert('Unauthorized. Please login.');
    window.location.href = "../index.html";
  } else {
    alert('Failed to add entry');
  }
}

async function deleteEntry(index) {
  const res = await fetch(`/api/entries/${index}`, { method: 'DELETE' });
  if (res.ok) {
    fetchEntries();
  } else if (res.status === 401) {
    alert('Unauthorized. Please login.');
    window.location.href = "../index.html";
  } else {
    alert('Failed to delete entry');
  }
}

function toggleVisibility(button) {
  const input = button.previousElementSibling;
  if (input.type === 'password') {
    input.type = 'text';
    button.textContent = '🙈';
  } else {
    input.type = 'password';
    button.textContent = '👁️';
  }
}

async function logout() {
  await fetch('/api/logout', { method: 'POST' });
  alert('Logged out');
  window.location.href = "../index.html";
}

fetchEntries();
