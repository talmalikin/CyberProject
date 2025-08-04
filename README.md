# CyberProject
✅ Application Overview
🎯 Functionality:
This is a Password Vault Web Application.
It allows users to:

Register with a username and password

Log in securely

Add records (Website, Email, Password)

View saved records (with hidden password and "show" button)

Optional future features:

Edit/Delete records

Cloud-based storage

Password encryption

Two-factor authentication

🧩 Components:
Component	Description
Frontend	Built with HTML, CSS, and vanilla JavaScript. Includes pages like login.html, register.html, and vault.html.
Backend	A server written in Go (Golang) that handles authentication and data storage through API endpoints.
Database	Currently implemented using an in-memory or file-based structure in Go (not a full database like PostgreSQL or MongoDB yet).
Static Files	Stored under the static/ folder (HTML/CSS/JS). Served by the Go backend.
APIs	RESTful API endpoints like:

POST /api/register

POST /api/login

POST /api/save

GET /api/load

🔁 Data Flow:
📥 Input:
User inputs through forms: username, password, and record details.

Submitted via JavaScript using fetch().

⚙️ Processing:
The Go server receives data through API calls.

It verifies credentials, saves data, or fetches data based on the route.

Basic logic: user authentication, data validation, storing records.

📤 Output:
If login is successful: user is redirected to vault.html.

Vault data is loaded and displayed.

If there's an error: error message is shown (e.g., “Invalid credentials”).

🔐 Sensitive Data Handling:
Data Type	How It's Handled
Passwords	Plain text or basic hashed (you should use proper hashing like bcrypt).
User sessions	Currently not implemented (recommended: use session tokens or JWTs).
Vault records	Stored per user; currently stored in memory or basic file (should be encrypted in production).
