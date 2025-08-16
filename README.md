# Pendaftaran Anggota HIMATIF - Backend

![Go](https://img.shields.io/badge/Go-00ADD8?style=for-the-badge&logo=go&logoColor=white)
![MongoDB](https://img.shields.io/badge/MongoDB-47A248?style=for-the-badge&logo=mongodb&logoColor=white)
![Paseto](https://img.shields.io/badge/Paseto-000000?style=for-the-badge&logo=paseto&logoColor=white)

Ini adalah layanan backend untuk aplikasi Pendaftaran Anggota Baru Himpunan Mahasiswa Teknik Informatika (HIMATIF). Dibangun menggunakan bahasa pemrograman Go dengan arsitektur yang bersih dan terpisah, layanan ini menyediakan REST API untuk mengelola data user, pendaftaran, dan otentikasi.

---

## 🚀 Fitur Utama

- **Otentikasi Aman**: Menggunakan **Paseto (PASETO)**, alternatif yang lebih aman dari JWT, untuk otentikasi berbasis token.
- **Manajemen User**: Pendaftaran (Register) dan Login user.
- **Manajemen Pendaftaran**: User dapat mengirimkan formulir pendaftaran lengkap dengan upload CV.
- **Dashboard Admin**: Admin dapat melihat, mengelola, dan mengubah status semua pendaftar.
- **Penyimpanan File**: CV yang diunggah disimpan di server dan path-nya direkam di database.
- **Arsitektur Decoupled**: Dirancang untuk bekerja secara terpisah dengan antarmuka frontend manapun.

---

## 🛠️ Teknologi yang Digunakan

- **Bahasa**: Go (v1.18+)
- **Database**: MongoDB Atlas
- **Router**: Chi (v5)
- **Otentikasi**: Paseto (v2)
- **Driver DB**: `go.mongodb.org/mongo-driver`
- **Lainnya**: `godotenv` untuk manajemen environment, `bcrypt` untuk hashing password.

---

## 📋 Prasyarat

Sebelum memulai, pastikan Anda telah menginstal:
- [Go](https://golang.org/doc/install) (versi 1.18 atau lebih baru)
- [MongoDB Atlas Account](https://www.mongodb.com/cloud/atlas) (Anda bisa menggunakan cluster gratis)
- Git

---

## ⚙️ Instalasi & Konfigurasi Lokal

1.  **Clone repository ini:**
    ```bash
    git clone [https://github.com/syalwa/pendaftaran-backend.git](https://github.com/syalwa/pendaftaran-backend.git)
    cd pendaftaran-backend
    ```

2.  **Siapkan file environment:**
    Buat file bernama `.env` di direktori utama dan isi dengan format berikut:
    ```env
    # Ambil dari MongoDB Atlas (klik Connect -> Drivers)
    MONGO_URI="mongodb+srv://<user>:<password>@<cluster-url>/<db-name>?retryWrites=true&w=majority"
    MONGO_DATABASE="himatif_db"

    # Kunci rahasia untuk Paseto (HARUS 32 karakter)
    PASETO_SECRET_KEY="R4nd0mS3cr3tK3yF0rP4s3t0Appl1c4t"

    # Port untuk server backend
    SERVER_PORT=":8080"
    ```

3.  **Instal dependensi:**
    Go akan mengunduh semua dependensi yang dibutuhkan secara otomatis. Anda bisa merapikannya dengan:
    ```bash
    go mod tidy
    ```

4.  **Jalankan server:**
    ```bash
    go run ./cmd/main/main.go
    ```
    Server akan berjalan di `http://localhost:8080`.

---

## 📝 Endpoint API

### Otentikasi
- `POST /register`: Mendaftarkan user baru.
- `POST /login`: Login user dan mendapatkan token Paseto.

### User (Memerlukan Token)
- `GET /api/user/profile`: Mendapatkan detail profil user yang sedang login.
- `POST /api/user/registration`: Mengirimkan formulir pendaftaran (termasuk upload CV).
- `GET /api/user/my-registration`: Mendapatkan status pendaftaran user yang sedang login.

### Admin (Memerlukan Token & Role Admin)
- `GET /api/admin/registrations-with-details`: Mendapatkan daftar semua pendaftar beserta detailnya.
- `PATCH /api/admin/registrations/{id}`: Memperbarui detail pendaftaran (status, jadwal wawancara).
- `GET /api/uploads/{nama_file_cv}`: Mengakses file CV yang telah diunggah.

---