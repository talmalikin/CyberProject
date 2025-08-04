# SECURITY.md — Security Documentation

## 1. Application Description and Components

### Main Functionality  
Web application for password management: users can register, login, save password entries (website, email, password), view and delete them.

### Components  
- **Frontend:** Static HTML/JS files served from `./static`.  
- **Backend:** Go HTTP server.  
- **Data Storage:** In-memory maps (no DB).  
- **User & Session Management:** Maps for users, sessions and password data.  
- **API:** RESTful endpoints for register, login, session management, saving and retrieving passwords.

## 2. Data Flow

| Step                 | Description                                                                               |
|----------------------|-------------------------------------------------------------------------------------------|
| User Registration    | Receives username and password; password hashed (bcrypt) and stored in `users` map.       |
| Login                | Checks password against stored hash, generates random sessionID and stores in `sessions`. |
| Save Password Entry  | Receives password entry, encrypts password field with AES-GCM, stores under `userData`.  |
| Retrieve Passwords   | Decrypts stored passwords and returns to client.                                         |
| Delete Password Entry| Deletes entry by index, after user authentication.                                       |

## 3. Potential Threats (Based on STRIDE)

| Threat Type          | Details                                          | Likelihood | Impact | Notes and Mitigations                   |
|----------------------|-------------------------------------------------|------------|--------|----------------------------------------|
| Spoofing             | Guessable SessionID → impersonation             | Medium     | High   | Strong random SessionID, HttpOnly cookie |
| Tampering            | Unauthorized data modification (Entries/Sessions)| Medium   | Medium | Mutex locks, user authentication on APIs |
| Repudiation (DoS)    | Flooding requests                                | Low        | High   | Currently no limits; recommend Rate Limiting |
| Information Disclosure| Exposure of passwords or sessions               | Medium     | High   | Passwords hashed/encrypted, HttpOnly cookies |
| Denial of Service    | Service unavailability (like DoS)                | Low        | High   | Advanced protections recommended        |
| Elevation of Privilege| Gaining unauthorized access                      | Low        | High   | Session verification, proper authorization |

## 4. Risk Assessment and Prioritization

| Threat                | Likelihood | Impact | Priority      |
|-----------------------|------------|--------|---------------|
| Session ID spoofing   | Medium     | High   | High          |
| Password disclosure   | Medium     | High   | High          |
| Data tampering        | Medium     | Medium | Medium        |
| DoS                   | Low        | High   | Medium        |
| Privilege escalation  | Low        | High   | High          |

## 5. Mitigations Implemented

- **Session Management:**  
  * Strong 32-byte random session IDs.  
  * HttpOnly cookies prevent JS access.  
  * Optionally Secure cookies with HTTPS.

- **Password Security:**  
  * bcrypt hashed user passwords.  
  * Secure password comparison.

- **Entry Password Encryption:**  
  * AES-GCM symmetric encryption of stored passwords.  
  * Encryption key hardcoded (recommend environment variable in production).  
  * Decryption before display.

- **Input Sanitization:**  
  * Trim spaces from inputs.

- **Concurrency Control:**  
  * Mutex locking for concurrent access.

- **Proper HTTP Methods:**  
  * Enforced POST/GET/DELETE methods.

## 6. Recommendations for Further Improvement

- Use external secure storage for encryption keys.  
- Add session expiration and cleanup.  
- Implement Rate Limiting for DoS protection.  
- Use HTTPS and Secure cookies.  
- Add RBAC authorization.  
- Enhanced error handling without sensitive info exposure.  
- Logging and anomaly detection.

---

