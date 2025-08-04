package main

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"
	"sync"
)

type Entry struct {
	Website  string `json:"website"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

var (
	users    = map[string]string{}       // username -> password
	sessions = map[string]string{}       // sessionID -> username
	userData = map[string][]Entry{}      // username -> entries
	mu       sync.Mutex
)

// helper: מקבל שם משתמש מ-cookie
func getUserFromRequest(r *http.Request) (string, bool) {
	cookie, err := r.Cookie("session_id")
	if err != nil {
		return "", false
	}
	mu.Lock()
	defer mu.Unlock()
	username, ok := sessions[cookie.Value]
	return username, ok
}

func main() {
	http.HandleFunc("/api/register", registerHandler)
	http.HandleFunc("/api/login", loginHandler)
	http.HandleFunc("/api/logout", logoutHandler)
	http.HandleFunc("/api/vault", vaultHandler)      // GET רשימות סיסמאות
	http.HandleFunc("/api/add", addHandler)          // POST הוספת סיסמה
	http.HandleFunc("/api/entries/", deleteHandler)  // DELETE מחיקת סיסמה לפי אינדקס

	// מאפשר שרת סטטי (למשל לשרת את קבצי frontend) אם רוצים
	http.Handle("/", http.FileServer(http.Dir("./static")))

	println("Server started at http://localhost:8080")
	http.ListenAndServe(":8080", nil)
}

func registerHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Only POST allowed", http.StatusMethodNotAllowed)
		return
	}
	type creds struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}
	var c creds
	if err := json.NewDecoder(r.Body).Decode(&c); err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}
	if c.Username == "" || c.Password == "" {
		http.Error(w, "Username and password required", http.StatusBadRequest)
		return
	}

	mu.Lock()
	defer mu.Unlock()
	if _, exists := users[c.Username]; exists {
		http.Error(w, "User exists", http.StatusBadRequest)
		return
	}

	users[c.Username] = c.Password
	userData[c.Username] = []Entry{}
	w.WriteHeader(http.StatusOK)
}

func loginHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Only POST allowed", http.StatusMethodNotAllowed)
		return
	}
	type creds struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}
	var c creds
	if err := json.NewDecoder(r.Body).Decode(&c); err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	mu.Lock()
	pass, ok := users[c.Username]
	mu.Unlock()
	if !ok || pass != c.Password {
		http.Error(w, "Invalid credentials", http.StatusUnauthorized)
		return
	}

	sessionID := "sess_" + c.Username
	mu.Lock()
	sessions[sessionID] = c.Username
	mu.Unlock()

	http.SetCookie(w, &http.Cookie{
		Name:  "session_id",
		Value: sessionID,
		Path:  "/",
	})

	w.WriteHeader(http.StatusOK)
}

func logoutHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Only POST allowed", http.StatusMethodNotAllowed)
		return
	}
	cookie, err := r.Cookie("session_id")
	if err == nil {
		mu.Lock()
		delete(sessions, cookie.Value)
		mu.Unlock()
		http.SetCookie(w, &http.Cookie{
			Name:   "session_id",
			Value:  "",
			Path:   "/",
			MaxAge: -1,
		})
	}
	w.WriteHeader(http.StatusOK)
}

func vaultHandler(w http.ResponseWriter, r *http.Request) {
	username, ok := getUserFromRequest(r)
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}
	if r.Method != http.MethodGet {
		http.Error(w, "Only GET allowed", http.StatusMethodNotAllowed)
		return
	}
	mu.Lock()
	entries := userData[username]
	mu.Unlock()

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(entries)
}

func addHandler(w http.ResponseWriter, r *http.Request) {
	username, ok := getUserFromRequest(r)
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}
	if r.Method != http.MethodPost {
		http.Error(w, "Only POST allowed", http.StatusMethodNotAllowed)
		return
	}
	var entry Entry
	if err := json.NewDecoder(r.Body).Decode(&entry); err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	mu.Lock()
	userData[username] = append(userData[username], entry)
	mu.Unlock()
	w.WriteHeader(http.StatusOK)
}

func deleteHandler(w http.ResponseWriter, r *http.Request) {
	username, ok := getUserFromRequest(r)
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}
	if r.Method != http.MethodDelete {
		http.Error(w, "Only DELETE allowed", http.StatusMethodNotAllowed)
		return
	}
	parts := strings.Split(r.URL.Path, "/")
	if len(parts) < 4 {
		http.Error(w, "Index missing", http.StatusBadRequest)
		return
	}
	indexStr := parts[3]
	index, err := strconv.Atoi(indexStr)
	if err != nil {
		http.Error(w, "Invalid index", http.StatusBadRequest)
		return
	}

	mu.Lock()
	defer mu.Unlock()
	entries := userData[username]
	if index < 0 || index >= len(entries) {
		http.Error(w, "Index out of range", http.StatusBadRequest)
		return
	}

	userData[username] = append(entries[:index], entries[index+1:]...)
	w.WriteHeader(http.StatusOK)
}
