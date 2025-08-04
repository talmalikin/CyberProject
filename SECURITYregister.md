# SECURITY.md — Register Page Security Documentation

## 1. Application Description and Components

**Main Functionality:**  
Frontend page for user registration: collecting username and password, sending to backend API `/api/register`.

**Components:**  
- HTML: register.html with username and password input fields  
- CSS: register.css for styling  
- JavaScript: register.js handling user input, validation, and communication with backend API

---

## 2. Data Flow

| Step                  | Description                                                                           |
|-----------------------|---------------------------------------------------------------------------------------|
| User Input            | User enters username and password in input fields                                     |
| Client Validation     | Username and password validated on client-side for format and length                  |
| Submit Registration   | Sends POST request with JSON payload `{ username, password }` to `/api/register`      |
| Handle Response       | Shows success or error message based on server response                              |
| Redirect on Success   | Redirects to login page (`index.html`) only if registration succeeded                 |

---

## 3. Potential Client-Side Threats

| Threat Type            | Details                                                                                   | Likelihood | Impact | Recommendations and Mitigations                      |
|------------------------|-------------------------------------------------------------------------------------------|------------|--------|-----------------------------------------------------|
| Invalid User Input      | Username or password fields empty or malformed (e.g., too short, invalid characters)      | High       | High   | Add client-side validation: username regex, min length, no empty fields |
| Use of alert()          | Alerts block UI, poor UX and possible information leakage                                 | Medium     | Medium | Replace `alert()` calls with a dedicated `<div id="message">` for messages |
| Poor Network Error Handling | Missing or insufficient handling of network or server errors                            | Medium     | Medium | Wrap fetch calls in `try...catch`, check `res.ok`, display proper messages |
| Leading/Trailing Spaces | User input with unwanted spaces causing subtle bugs                                      | Medium     | Medium | Use `.trim()` on inputs before validation and submission |
| Premature Redirection  | Redirect to login page even on registration failure or server error                      | Medium     | Medium | Redirect only when `res.ok === true` (already implemented) |

---

## 4. Risk Assessment and Prioritization

| Threat                 | Likelihood | Impact | Priority  |
|------------------------|------------|--------|-----------|
| Invalid user input     | High       | High   | High      |
| Use of alert()         | Medium     | Medium | Medium    |
| Network error handling | Medium     | Medium | Medium    |
| Input trimming issues  | Medium     | Medium | Medium    |
| Premature redirect     | Medium     | Medium | Medium    |

---

## 5. Implemented Mitigations and Best Practices

| Mitigation                             | Description                                                      |
|---------------------------------------|------------------------------------------------------------------|
| Input trimming                        | `.trim()` used to remove whitespace before validation and submit |
| Client-side input validation          | Username regex to allow only letters, numbers, underscore; min length 3; password min length 8 |
| Replacement of alert()                 | Use of a `<div id="message">` to show user-friendly messages     |
| Network error handling                 | Use of `try...catch` around `fetch`, checking `res.ok`            |
| Conditional redirect                   | Redirect to login page only on successful registration (`res.ok`) |

---

## 6. Recommendations for Future Improvement

| Recommendation                        | Description                                                      |
|-------------------------------------|------------------------------------------------------------------|
| Strengthen password rules            | Enforce stronger password complexity (special chars, uppercase) |
| Add CAPTCHA or rate limiting         | Prevent automated registration abuse                            |
| Use HTTPS                           | Ensure secure transmission of credentials                        |
| Server-side validation and sanitization | Never rely solely on client-side validation                      |
| Use Content Security Policy (CSP)    | Prevent XSS and injection attacks                                |
| Logging and monitoring                | Track registration attempts and failures                        |

---

**Last updated:** August 4, 2025
