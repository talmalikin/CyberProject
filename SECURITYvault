# SECURITY.md — Vault Frontend Security Documentation

## 1. Application Description and Components

### Main Functionality  
Frontend page for managing stored password entries: displaying, adding, and deleting user passwords.

### Components  
- **HTML:** vault.html with input fields and layout.  
- **CSS:** vault.css for styling.  
- **JavaScript:** vault.js handling user interaction and communication with backend APIs.

## 2. Data Flow

| Step                   | Description                                                                                   |
|------------------------|-----------------------------------------------------------------------------------------------|
| Fetch Password Entries  | On page load, fetches `/api/vault` to retrieve encrypted entries after user authentication.   |
| Display Entries        | Dynamically renders entries including website, email, and masked password fields.             |
| Add Entry              | Sends new entry data to `/api/add` via POST after client-side input validation.                |
| Delete Entry           | Sends DELETE request to `/api/entries/{index}` to remove an entry by index.                    |
| Logout                 | Sends POST request to `/api/logout` and redirects to login page.                              |

## 3. Potential Client-Side Threats

| Threat Type            | Details                                                                                            | Likelihood | Impact | Recommendations and Mitigations                              |
|------------------------|--------------------------------------------------------------------------------------------------|------------|--------|-------------------------------------------------------------|
| Cross-Site Scripting (XSS) | Usage of `innerHTML` to insert user data (website, email, password) without proper sanitization. | Medium     | High   | Replace `innerHTML` with safe DOM manipulation (`createElement`, `textContent`). |
| Lack of Input Validation | No checks if website, email, or password fields are valid before sending to the server.           | High       | Medium | Implement client-side validation: website format, email syntax, password length. |
| Use of alert()          | Alerts block UI and degrade user experience.                                                     | Medium     | Medium | Replace alerts with dynamic messages inside dedicated page elements.           |
| Insufficient Network Error Handling | Fetch calls lack comprehensive error handling (e.g., network disconnects).                | Medium     | Medium | Wrap fetch calls in try/catch blocks and display meaningful error messages.    |
| Password Exposure Time  | Password visible in plaintext while toggled visible, no auto-hide or copy prevention.            | Medium     | Medium | Auto-hide password after timeout, prevent copying visible passwords.            |

## 4. Risk Assessment and Prioritization

| Threat                | Likelihood | Impact | Priority  |
|-----------------------|------------|--------|-----------|
| XSS via innerHTML     | Medium     | High   | High      |
| Lack of Input Validation | High     | Medium | High      |
| Use of alert()        | Medium     | Medium | Medium    |
| Network Error Handling | Medium    | Medium | Medium    |
| Password Exposure Time | Medium    | Medium | Medium    |

## 5. Implemented Mitigations and Best Practices

- Trim whitespace from user inputs before submission.  
- Redirect unauthorized users (401) to the login page.  
- Toggle password visibility securely with a button.  
- Clear input fields after successful add operations.  
- Use `async/await` and `try/catch` for proper fetch call handling.

## 6. Recommendations for Future Improvement

- Replace all `innerHTML` usage with safe DOM element creation (`document.createElement`) and set `textContent`.  
- Add comprehensive client-side validation for website, email, and password fields.  
- Replace `alert()` calls with dynamic in-page message displays.  
- Wrap all fetch calls in `try/catch` with user-friendly error reporting.  
- Add password display security: auto-hide passwords after a timeout, disable copy/paste when visible.  
- Implement Content Security Policy (CSP) headers server-side to reduce XSS risk.  
- Use HTTPS and Secure, HttpOnly cookies on the server side.  
- Implement client-side error logging and monitoring.

---

**Last updated:** August 4, 2025
