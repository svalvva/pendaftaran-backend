# Pendaftaran Anggota HIMATIF - Backend 🚀

![Go](https://img.shields.io/badge/Go-00ADD8?style=for-the-badge&logo=go&logoColor=white)
![MongoDB](https://img.shields.io/badge/MongoDB-47A248?style=for-the-badge&logo=mongodb&logoColor=white)
![Cloudinary](https://img.shields.io/badge/Cloudinary-3448C5?style=for-the-badge&logo=Cloudinary&logoColor=white)
![Paseto](https://img.shields.io/badge/Paseto-000000?style=for-the-badge&logo=paseto&logoColor=white)

Selamat datang di layanan backend untuk Aplikasi Pendaftaran Anggota Baru HIMATIF! Dibuat dengan Go, layanan ini menyediakan REST API yang tangguh dan aman untuk mengelola seluruh alur pendaftaran, mulai dari otentikasi pengguna hingga manajemen file di *cloud*.

---

## ✨ Fitur Unggulan

-   **Otentikasi Modern & Aman**: Menggunakan **Paseto (PASETO)**, alternatif superior dari JWT, untuk token yang lebih aman.
-   **Manajemen User Lengkap**: Registrasi dan Login untuk calon anggota baru.
-   **Pendaftaran Online**: Pengguna dapat mengisi formulir pendaftaran secara lengkap.
-   **Cloud File Uploads**: CV dan sertifikat diunggah langsung ke **Cloudinary**, memastikan penyimpanan file yang efisien dan aman.
-   **Dashboard Admin Komprehensif**: Admin dapat melihat, mengelola, memperbarui status (termasuk aksi massal), dan menghapus pendaftar.
-   **Manajemen Informasi**: Admin dapat membuat, membaca, memperbarui, dan menghapus pengumuman atau informasi untuk semua pengguna.

---

## 🛠️ Tumpukan Teknologi

-   **Bahasa**: Go (v1.23+)
-   **Database**: MongoDB Atlas
-   **Router**: Chi (v5)
-   **Otentikasi**: Paseto (v2)
-   **Penyimpanan File**: Cloudinary
-   **Driver DB**: `go.mongodb.org/mongo-driver`
-   **Lainnya**: `godotenv` untuk manajemen *environment*, `bcrypt` untuk *hashing* password.

---

## ⚙️ Instalasi & Konfigurasi Lokal

1.  **Clone repository ini:**
    ```bash
    git clone [https://github.com/svalvva/pendaftaran-backend.git](https://github.com/svalvva/pendaftaran-backend.git)
    cd pendaftaran-backend
    ```

2.  **Siapkan file environment:**
    Buat file `.env` di direktori utama dan isi dengan format berikut:
    ```env
    # Ambil dari MongoDB Atlas (klik Connect -> Drivers)
    MONGO_URI="mongodb+srv://<user>:<password>@<cluster-url>/<db-name>?retryWrites=true&w=majority"
    MONGO_DATABASE="himatif_db"

    # Kunci rahasia untuk Paseto (HARUS 32 karakter)
    PASETO_SECRET_KEY="R4nd0mS3cr3tK3yF0rP4s3t0Appl1c4t"

    # Port untuk server backend
    SERVER_PORT=":8080"
    
    # Kredensial dari akun Cloudinary Anda
    CLOUDINARY_CLOUD_NAME="<your_cloud_name>"
    CLOUDINARY_API_KEY="<your_api_key>"
    CLOUDINARY_API_SECRET="<your_api_secret>"
    ```

3.  **Instal dependensi:**
    ```bash
    go mod tidy
    ```

4.  **Jalankan server:**
    ```bash
    go run ./main.go
    ```
    Server akan berjalan di `http://localhost:8080`.

---

## 📝 Endpoint API

### Otentikasi
- `POST /register`: Mendaftarkan user baru.
- `POST /login`: Login user dan mendapatkan token Paseto.

### Pengguna (Memerlukan Token)
- `GET /api/user/profile`: Mendapatkan detail profil user yang sedang login.
- `POST /api/user/registration`: Mengirimkan formulir pendaftaran (termasuk upload CV & sertifikat).
- `GET /api/user/my-registration`: Mendapatkan status pendaftaran user yang sedang login.
- `GET /api/info`: Mendapatkan semua informasi/pengumuman terbaru.

### Admin (Memerlukan Token & Role Admin)
- `GET /api/admin/registrations-with-details`: Mendapatkan daftar semua pendaftar beserta detailnya.
- `GET /api/admin/users`: Mendapatkan daftar semua pengguna terdaftar.
- `PATCH /api/admin/registrations/{id}`: Memperbarui detail pendaftaran (status, jadwal wawancara, dll).
- `PATCH /api/admin/registrations/bulk-update`: Memperbarui status beberapa pendaftar sekaligus.
- `DELETE /api/admin/registrations/{id}`: Menghapus data pendaftaran.
- `POST /api/admin/info`: Membuat informasi/pengumuman baru.
- `PUT /api/admin/info/{id}`: Memperbarui informasi yang sudah ada.
- `DELETE /api/admin/info/{id}`: Menghapus informasi.

---
