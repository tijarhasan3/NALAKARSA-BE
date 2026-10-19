# Update FE: Verifikasi Email dan Reset Password

## Status Fitur

Backend saat ini mendukung:

- Email verifikasi akun dengan kode 6 digit.
- Pengiriman ulang kode verifikasi.
- Email selamat datang setelah email berhasil diverifikasi.
- Email reset password berbentuk HTML dengan tombol **Atur Ulang Password**.
- Email notifikasi setelah password berhasil diubah.

Semua endpoint menggunakan prefix:

```text
/api/v1
```

## Alur Setelah Register

1. FE memanggil `POST /auth/register`.
2. Backend membuat akun dan mengirim kode verifikasi ke email user.
3. FE mengarahkan user ke halaman verifikasi email.
4. User memasukkan kode 6 digit dari email.
5. FE memanggil `POST /auth/verify-email`.
6. Setelah berhasil, backend mengirim email selamat datang.

Kode verifikasi berlaku selama 15 menit dan tidak dikembalikan dalam response register.

## Verifikasi Email Setelah Registrasi

Setelah registrasi berhasil, backend mengirim email HTML dengan kode verifikasi 6 digit ke email user. Kode berlaku selama 15 menit.

### Verifikasi kode

```http
POST /api/v1/auth/verify-email
Content-Type: application/json
```

```json
{
  "email": "user@example.com",
  "code": "123456"
}
```

Response sukses:

```json
{
  "message": "Email berhasil diverifikasi",
  "data": null
}
```

### Kirim ulang kode

```http
POST /api/v1/auth/resend-verification-email
Content-Type: application/json
```

```json
{
  "email": "user@example.com"
}
```

Setelah kode berhasil diverifikasi, backend mengirim email selamat datang.

Jika kode salah atau kedaluwarsa, user dapat meminta kode baru melalui endpoint resend.

## Ringkasan

Fitur lupa password sekarang mengirimkan link reset password melalui email. Token reset tidak lagi dikembalikan pada response endpoint `forgot-password`.

## Alur FE

1. User membuka halaman **Lupa Password**.
2. User mengisi email, security question, dan security answer.
3. FE memanggil endpoint `POST /api/v1/auth/forgot-password`.
4. Tampilkan pesan bahwa link reset telah dikirim melalui email.
5. User menekan tombol **Atur Ulang Password** dari email:

   ```text
   /reset-password?token=RESET_TOKEN
   ```

6. Halaman reset password mengambil token dari query parameter `token`.
7. User mengisi password baru dan FE memanggil endpoint `POST /api/v1/auth/reset-password`.
8. Setelah berhasil, arahkan user ke halaman login.

## Endpoint: Request Reset Password

### Request

```http
POST /api/v1/auth/forgot-password
Content-Type: application/json
```

```json
{
  "email": "user@example.com",
  "security_question": "favorite_animal",
  "security_answer": "kucing"
}
```

Nilai `security_question` yang valid:

- `favorite_animal`
- `mother_name`

### Success response

HTTP `200 OK`:

```json
{
  "message": "Jika data valid, link reset password telah dikirim ke email Anda",
  "data": null
}
```

Response ini tidak mengandung `reset_token`. Token hanya berada di link tombol email.

## Endpoint: Reset Password

### Request

```http
POST /api/v1/auth/reset-password
Content-Type: application/json
```

```json
{
  "token": "TOKEN_DARI_QUERY_PARAMETER",
  "new_password": "PasswordBaru123!"
}
```

### Success response

HTTP `200 OK`:

```json
{
  "message": "Password reset successfully",
  "data": null
}
```

Setelah berhasil, semua sesi login lama user akan dicabut dan user harus login kembali.

## Validasi password

Password baru harus memenuhi ketentuan berikut:

- Minimal 6 karakter
- Maksimal 72 karakter
- Mengandung huruf besar
- Mengandung huruf kecil
- Mengandung angka
- Mengandung karakter khusus

Contoh valid:

```text
PasswordBaru123!
```

## Penanganan error FE

Untuk error verifikasi email, tampilkan pesan yang mudah dipahami:

```text
Kode verifikasi tidak valid atau sudah kedaluwarsa.
```

Untuk error pengiriman email:

```text
Email belum dapat dikirim. Silakan coba lagi beberapa saat kemudian.
```

Jika token invalid atau kedaluwarsa, tampilkan pesan seperti:

```text
Link reset password tidak valid atau sudah kedaluwarsa. Silakan minta link baru.
```

Jika email reset gagal dikirim, tampilkan pesan umum dan minta user mencoba lagi. Jangan menampilkan SMTP error kepada user.

Jangan menampilkan kode verifikasi atau token reset di log browser maupun UI selain pada field input verifikasi.

## Contoh implementasi sederhana

Mengambil token dari URL:

```javascript
const params = new URLSearchParams(window.location.search);
const token = params.get("token");
```

Mengirim password baru:

```javascript
await fetch("/api/v1/auth/reset-password", {
  method: "POST",
  headers: {
    "Content-Type": "application/json"
  },
  body: JSON.stringify({
    token,
    new_password: newPassword
  })
});
```

## Catatan konfigurasi

URL link pada email menggunakan nilai backend `FRONTEND_URL`. Untuk development:

```env
FRONTEND_URL=http://localhost:5173
```

Untuk production, nilai tersebut harus diarahkan ke domain frontend production, misalnya:

```env
FRONTEND_URL=https://app.nalakarsa.id
```
