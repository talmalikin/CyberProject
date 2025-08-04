package main

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"io"
	"net/http"
	"strconv"
	"strings"
	"sync"

	"golang.org/x/crypto/bcrypt"
)

type Entry struct {
	Website  string `json:"website"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

var (
	users    = map[string][]byte{}  // username -> hashed password
	sessions = map[string]string{}  // sessionID -> username
	userData = map[string][]Entry{} // username -> entries
	mu       sync.Mutex

	// מפתח הצפנה (32 בתים = AES-256)
	encryptionKey = []byte("01234567890123456789012345678901") // החלף למפתח סודי באמת
)

// --- פונקציות הצפנה ---

func encrypt(plaintext string) (string, error) {
	block, err := aes.NewCipher(encryptionKey)
	if err != nil {
		return "", err
	}
	aesGCM, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}
	nonce := make([]byte, aesGCM.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return "", err
	}
	ciphertext := aesGCM.Seal(nonce, nonce, []byte(plaintext), nil)
	return base64.StdEncoding.EncodeToString(ciphertext), nil
}

func decrypt(ciphertextEnc string) (string, error) {
	ciphertext, err := base64.StdEncoding.DecodeString(ciphertextEnc)
	if err != nil {
		return "", err
	}
	block, err := aes.NewCipher(encryptionKey)
	if err != nil {
		return "", err
	}
	aesGCM, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}
	nonceSize := aesGCM.NonceSize()
	if len(ciphertext) < nonceSize {
		return "", err
	}
	nonce, ciphertext := ciphertext[:nonceSize], ciphertext[nonceSize:]
	plaintext, err := aesGCM.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return "", err
	}
	return string(plaintext), nil
}

// --- ניהול sessions ---

func generateSessionID() (string, error) {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(b), nil
}

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

// --- handlers ---

func main() {
	http.HandleFunc("/api/register", registerHandler)
	http.HandleFunc("/api/login", loginHandler)
	http.HandleFunc("/api/logout", logoutHandler)
	http.HandleFunc("/api/vault", vaultHandler)
	http.HandleFunc("/api/add", addHandler)
	http.HandleFunc("/api/entries/", deleteHandler)

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

	hashedPass, err := bcrypt.GenerateFromPassword([]byte(c.Password), bcrypt.DefaultCost)
	if err != nil {
		http.Error(w, "Server error", http.StatusInternalServerError)
		return
	}

	mu.Lock()
	defer mu.Unlock()
	if _, exists := users[c.Username]; exists {
		http.Error(w, "User exists", http.StatusBadRequest)
		return
	}

	users[c.Username] = hashedPass
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
	hashedPass, ok := users[c.Username]
	mu.Unlock()
	if !ok {
		http.Error(w, "Invalid credentials", http.StatusUnauthorized)
		return
	}
	if err := bcrypt.CompareHashAndPassword(hashedPass, []byte(c.Password)); err != nil {
		http.Error(w, "Invalid credentials", http.StatusUnauthorized)
		return
	}

	sessionID, err := generateSessionID()
	if err != nil {
		http.Error(w, "Server error", http.StatusInternalServerError)
		return
	}

	mu.Lock()
	sessions[sessionID] = c.Username
	mu.Unlock()

	http.SetCookie(w, &http.Cookie{
		Name:     "session_id",
		Value:    sessionID,
		Path:     "/",
		HttpOnly: true,
		// Secure: true, // הפעל בייצור עם HTTPS
		// SameSite: http.SameSiteStrictMode,
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

	decryptedEntries := make([]Entry, len(entries))
	for i, e := range entries {
		pass, err := decrypt(e.Password)
		if err != nil {
			pass = "ERROR_DECRYPTING"
		}
		decryptedEntries[i] = Entry{
			Website:  e.Website,
			Email:    e.Email,
			Password: pass,
		}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(decryptedEntries)
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

	entry.Website = strings.TrimSpace(entry.Website)
	entry.Email = strings.TrimSpace(entry.Email)
	entry.Password = strings.TrimSpace(entry.Password)

	encryptedPass, err := encrypt(entry.Password)
	if err != nil {
		http.Error(w, "Encryption error", http.StatusInternalServerError)
		return
	}
	entry.Password = encryptedPass

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
