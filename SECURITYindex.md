# SECURITY.md — Login Page Security Documentation

## 1. Application Description and Components

**Main Functionality:**  
Frontend login page where users input username and password, which are sent to the backend API `/api/login`. On success, users are redirected to the password vault page.

**Components:**  
- HTML (`index.html`)  
- CSS (`index.css`)  
- JavaScript (`index.js`)

---

## 2. Data Flow Overview

| Step              | Description                                           |
|-------------------|-------------------------------------------------------|
| User Input        | User enters username and password                      |
| Client Validation | Input validated for format and length                  |
| Submit Login      | Sends POST request with JSON `{username, password}`   |
| Handle Response   | Shows success or error messages                         |
| Redirect          | Redirects to `vault.html` only if login succeeded      |

---

## 3. Identified Client-Side Security Issues and Solutions

| Issue                               | Details                                                             | Likelihood | Impact | Solution                                                                                  |
|-----------------------------------|---------------------------------------------------------------------|------------|--------|-------------------------------------------------------------------------------------------|
| Password input shown as plain text | Password `<input>` of type `text` exposes password characters       | Medium     | High   | Change to `<input type="password">` to mask user input                                  |
| Missing input validation           | Username/password not validated, allowing invalid or malicious data | High       | High   | Add regex validation for username (letters, numbers, underscore), minimum password length |
| Use of `alert()` for messages     | Alerts are intrusive and can be confusing or spoofed                | Medium     | Medium | Replace `alert()` with a dedicated message `<div id="message">` for safe user feedback    |
| Lack of error handling on network  | No `try...catch` or response status check may cause failures        | Medium     | Medium | Wrap fetch calls in `try...catch`; verify `res.ok` and show appropriate error messages     |
| Leading/trailing spaces in inputs  | Spaces can bypass validations or cause unexpected behavior          | Medium     | High   | Use `.trim()` on inputs before validation and submission                                 |
| Unconditional redirect             | Redirects even if login failed                                      | Medium     | Medium | Redirect only on successful login (`res.ok === true`)                                    |
| Potential client-side injection    | Unvalidated input could allow code injection if displayed insecurely| Low        | High   | Validate inputs; display feedback as text, avoid inserting raw HTML                      |

---

## 4. Risk Assessment

| Threat                        | Likelihood | Impact | Priority |
|------------------------------|------------|--------|----------|
| Invalid input & injection    | High       | High   | High     |
| Password exposure            | Medium     | High   | High     |
| Network error mishandling    | Medium     | Medium | Medium   |
| Intrusive alert usage        | Medium     | Medium | Medium   |
| Premature redirect           | Medium     | Medium | Medium   |

---

## 5. Mitigations Implemented

- Password `<input>` changed to type `password` to hide characters  
- Client-side validation: username regex, password length check, no empty fields  
- Non-intrusive user feedback via a dedicated `<div id="message">` instead of `alert()`  
- Proper network error handling with `try...catch` and `res.ok` checks  
- Input values trimmed of whitespace before processing  
- Redirect only occurs when login is successful

---

## 6. Recommendations for Further Improvement

- Enforce HTTPS for all connections  
- Add server-side validation and sanitization  
- Implement rate limiting and brute force protections  
- Add Content Security Policy (CSP) headers  
- Use HTTP-only, Secure cookies for session management  
- Consider multi-factor authentication (MFA)

---

_Last updated: 2025-08-04_
