# Test Bruno: Auth dan Email

## Environment

Buat environment Bruno dengan variabel:

```text
baseUrl = http://localhost:8080/api/v1
testEmail = email-kamu@example.com
```

Pastikan backend sudah berjalan dan `.env` Brevo sudah terisi.

## 1. Register

```http
POST {{baseUrl}}/auth/register
Content-Type: application/json
```

```json
{
  "email": "{{testEmail}}",
  "password": "Password123!",
  "security_question": "favorite_animal",
  "security_answer": "kucing",
  "role": "profesional",
  "firstName": "User",
  "lastName": "Test",
  "fullName": "User Test",
  "affiliation": "Nalakarsa",
  "location": "Jakarta",
  "expertise": "Technology",
  "industry": "Technology"
}
```

Expected result: HTTP `201` dan email kode verifikasi masuk ke inbox.

## 2. Verify email

Ambil kode 6 digit dari email, lalu jalankan:

```http
POST {{baseUrl}}/auth/verify-email
Content-Type: application/json
```

```json
{
  "email": "{{testEmail}}",
  "code": "123456"
}
```

Expected result: HTTP `200` dengan pesan `Email berhasil diverifikasi`.

## 3. Resend verification code

```http
POST {{baseUrl}}/auth/resend-verification-email
Content-Type: application/json
```

```json
{
  "email": "{{testEmail}}"
}
```

## 4. Request password reset

```http
POST {{baseUrl}}/auth/forgot-password
Content-Type: application/json
```

```json
{
  "email": "{{testEmail}}",
  "security_question": "favorite_animal",
  "security_answer": "kucing"
}
```

Ambil link reset dari email. Token adalah nilai setelah `token=` pada URL.

## 5. Reset password

```http
POST {{baseUrl}}/auth/reset-password
Content-Type: application/json
```

```json
{
  "token": "TOKEN_DARI_EMAIL",
  "new_password": "PasswordBaru123!"
}
```

Expected result: HTTP `200`, email notifikasi perubahan password masuk, dan refresh token lama dicabut.

## Catatan

- Kode verifikasi berlaku 15 menit.
- Token reset password berlaku 15 menit.
- Pastikan `SMTP_FROM` sudah diverifikasi di Brevo.
- Jika muncul `unauthorized IP`, tambahkan IP server pada pengaturan SMTP Brevo.
