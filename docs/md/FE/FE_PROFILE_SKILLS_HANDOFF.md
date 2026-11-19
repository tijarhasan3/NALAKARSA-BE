# Profile Skills Handoff untuk FE

## Ringkasan Perubahan

Sekarang terdapat dua konsep yang berbeda:

- `expertise`: bidang keahlian utama yang diisi saat registrasi. Nilainya tetap berada di profil user sebagai satu nilai.
- `skills`: kemampuan tambahan yang dapat dimiliki user dalam jumlah banyak dan dikelola dari halaman profile.

`expertise` tidak diganti oleh `skills`.

## Struktur Data

Relasi skill menggunakan tiga tabel:

```text
users
  id

skills
  id
  name

user_skills
  user_id
  skill_id
```

Satu user dapat memiliki banyak skill. Skill yang sama juga dapat digunakan oleh banyak user.

## Base URL

```text
{{url}}/api/v1
```

Gunakan header berikut untuk endpoint yang membutuhkan login:

```http
Authorization: Bearer {{token}}
Content-Type: application/json
```

## 1. Mencari Skill

Digunakan untuk mengisi combo box atau autocomplete.

```http
GET {{url}}/api/v1/skills/search?q=Machine&limit=20
```

Parameter:

| Parameter | Wajib | Keterangan |
|---|---:|---|
| `q` | Tidak | Kata kunci pencarian |
| `limit` | Tidak | Jumlah data, default 20, maksimum 50 |

Contoh response:

```json
{
  "data": [
    {
      "id": "uuid-skill",
      "name": "Machine Learning"
    }
  ],
  "message": "Skills retrieved successfully",
  "timestamp": "2026-10-07T00:00:00Z"
}
```

## 2. Mengambil Skill Profile Sendiri

```http
GET {{url}}/api/v1/users/me/skills
```

Wajib menggunakan token login.

Contoh response:

```json
{
  "data": [
    {
      "id": "uuid-skill-1",
      "name": "Python"
    },
    {
      "id": "uuid-skill-2",
      "name": "Data Analysis"
    }
  ],
  "message": "User skills retrieved successfully",
  "timestamp": "2026-10-07T00:00:00Z"
}
```

Jika user belum memiliki skill, `data` dikembalikan sebagai array kosong.

## 3. Menambahkan Skill yang Sudah Ada

```http
POST {{url}}/api/v1/users/me/skills
```

Request body:

```json
{
  "skill_id": "uuid-skill"
}
```

Jika skill sudah dimiliki user, backend tidak membuat relasi duplikat.

## 4. Menambahkan Skill Baru

Jika hasil pencarian kosong dan user memasukkan skill baru, gunakan endpoint yang sama.

```http
POST {{url}}/api/v1/users/me/skills
```

Request body:

```json
{
  "name": "Quantum Computing"
}
```

Backend akan:

1. Menormalisasi nama skill.
2. Mengecek skill dengan nama yang sama tanpa membedakan huruf besar dan kecil.
3. Membuat skill baru jika belum tersedia.
4. Menghubungkan skill ke profile user.

FE cukup memanggil satu endpoint.

## 5. Menghapus Skill dari Profile

```http
DELETE {{url}}/api/v1/users/me/skills/{skill_id}
```

Contoh:

```http
DELETE {{url}}/api/v1/users/me/skills/uuid-skill
```

Wajib menggunakan token login.

Penghapusan hanya menghapus hubungan user dengan skill. Data master pada tabel `skills` tidak ikut dihapus.

## 6. Skill pada Profile User Lain

```http
GET {{url}}/api/v1/skills/users/{user_id}
```

Endpoint ini dapat digunakan untuk menampilkan skill pada profile publik.

## Alur Combo Box Profile

1. Saat halaman profile dibuka, panggil `GET /users/me/skills`.
2. Saat user mengetik, panggil `GET /skills/search?q=...`.
3. Jika user memilih hasil yang tersedia, kirim `skill_id` ke `POST /users/me/skills`.
4. Jika user membuat skill baru, kirim `name` ke endpoint yang sama.
5. Setelah berhasil, tambahkan response ke daftar skill di UI.
6. Saat user menghapus skill, panggil endpoint `DELETE`.

## Catatan untuk FE

- `expertise` tetap digunakan pada form registrasi.
- `skills` dikelola pada halaman profile.
- Jangan mengirim `skill_id` dan `name` sekaligus. Gunakan salah satu.
- Simpan `id` skill dari response untuk kebutuhan hapus.
- Gunakan `data: []` sebagai kondisi profile tanpa skill.
- Endpoint `users/me/skills` membutuhkan token JWT.
- Skill baru tidak perlu dibuat melalui endpoint katalog secara terpisah.

## OpenAPI dan Bruno

Dokumentasi OpenAPI sudah diperbarui di:

```text
docs/openapi.yaml
```

Import ulang file tersebut ke Bruno agar folder endpoint `Profile Skills` muncul otomatis.
