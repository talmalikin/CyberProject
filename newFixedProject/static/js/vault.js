const messageDiv = document.getElementById('message');

document.getElementById('toggle-add-btn').addEventListener('click', toggleAddFields);
document.querySelector('.hide-btn').addEventListener('click', toggleAddFields);
document.getElementById('add-entry-btn').addEventListener('click', addEntry);
document.querySelector('.logout-btn').addEventListener('click', logout);

function showMessage(msg, isError = false) {
  messageDiv.textContent = msg;
  messageDiv.style.color = isError ? 'red' : 'yellow';
  setTimeout(() => {
    messageDiv.textContent = '';
  }, 4000);
}

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

function escapeHtml(text) {
  return text.replace(/[&<>"']/g, function(m) {
    return ({
      '&': '&amp;',
      '<': '&lt;',
      '>': '&gt;',
      '"': '&quot;',
      "'": '&#39;'
    })[m];
  });
}

function createEntryElement(item, index) {
  const entryDiv = document.createElement('div');
  entryDiv.className = 'entry';

  const leftDiv = document.createElement('div');
  leftDiv.className = 'entry-left';

  const iconBox = document.createElement('div');
  iconBox.className = 'icon-box';
  iconBox.textContent = item.website.slice(0, 2).toUpperCase();

  const infoDiv = document.createElement('div');
  infoDiv.className = 'entry-info';

  const websiteDiv = document.createElement('div');
  websiteDiv.className = 'entry-line';
  websiteDiv.innerHTML = `<strong>${escapeHtml(item.website)}</strong>`;

  const emailDiv = document.createElement('div');
  emailDiv.className = 'entry-line';
  emailDiv.innerHTML = `<span class="icon">👤</span> ${escapeHtml(item.email)}`;

  const passwordDiv = document.createElement('div');
  passwordDiv.className = 'entry-line password-box';

  const input = document.createElement('input');
  input.className = 'password-input';
  input.type = 'password';
  input.value = item.password;
  input.readOnly = true;

  const eyeBtn = document.createElement('button');
  eyeBtn.className = 'eye-btn';
  eyeBtn.textContent = '👁️';
  eyeBtn.addEventListener('click', () => toggleVisibility(eyeBtn));

  passwordDiv.appendChild(document.createTextNode('🔒 '));
  passwordDiv.appendChild(input);
  passwordDiv.appendChild(eyeBtn);

  infoDiv.appendChild(websiteDiv);
  infoDiv.appendChild(emailDiv);
  infoDiv.appendChild(passwordDiv);

  leftDiv.appendChild(iconBox);
  leftDiv.appendChild(infoDiv);

  const deleteBtn = document.createElement('span');
  deleteBtn.className = 'delete-btn';
  deleteBtn.textContent = '🗑️';
  deleteBtn.addEventListener('click', () => deleteEntry(index));

  entryDiv.appendChild(leftDiv);
  entryDiv.appendChild(deleteBtn);

  return entryDiv;
}

async function fetchEntries() {
  try {
    const res = await fetch('/api/vault');
    if (res.status === 401) {
      showMessage('Unauthorized. Please login.', true);
      window.location.href = "../index.html";
      return;
    }
    if (!res.ok) {
      showMessage('Failed to load entries', true);
      return;
    }
    const data = await res.json();
    const container = document.getElementById('entries');
    container.innerHTML = '';
    data.forEach((item, index) => {
      const entryEl = createEntryElement(item, index);
      container.appendChild(entryEl);
    });
  } catch (err) {
    showMessage('Network error', true);
  }
}

async function addEntry() {
  const website = document.getElementById('website').value.trim();
  const email = document.getElementById('email').value.trim();
  const password = document.getElementById('password').value.trim();

  if (!website || !email || !password) {
    showMessage('Please fill all fields.', true);
    return;
  }

  // פשוט ולידציה בסיסית (אפשר להרחיב)
  
  if (!email.includes('@')) {
    showMessage('Invalid email address.', true);
    return;
  }
  

  try {
    const res = await fetch('/api/add', {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({ website, email, password })
    });

    if (res.ok) {
      clearAddFields();
      toggleAddFields();
      fetchEntries();
      showMessage('Password added successfully.');
    } else if (res.status === 401) {
      showMessage('Unauthorized. Please login.', true);
      window.location.href = "../index.html";
    } else {
      showMessage('Failed to add entry.', true);
    }
  } catch (err) {
    showMessage('Network error.', true);
  }
}

async function deleteEntry(index) {
  try {
    const res = await fetch(`/api/entries/${index}`, { method: 'DELETE' });
    if (res.ok) {
      fetchEntries();
      showMessage('Entry deleted.');
    } else if (res.status === 401) {
      showMessage('Unauthorized. Please login.', true);
      window.location.href = "../index.html";
    } else {
      showMessage('Failed to delete entry.', true);
    }
  } catch {
    showMessage('Network error.', true);
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
  try {
    await fetch('/api/logout', { method: 'POST' });
    showMessage('Logged out.');
    window.location.href = "../index.html";
  } catch {
    showMessage('Network error during logout.', true);
  }
}

// קריאה ראשונית לטעינת הרשומות
fetchEntries();
