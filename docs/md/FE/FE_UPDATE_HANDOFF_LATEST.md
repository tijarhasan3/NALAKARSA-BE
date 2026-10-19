# FE Handoff: Update Terbaru Nalakarsa

Dokumen ini adalah ringkasan kontrak FE terbaru untuk reference data, registrasi, dan update profile.

## Base URL

```text
url = https://DOMAIN-RAILWAY/api/v1
```

Contoh:

```http
GET {{url}}/knowledge/search?q=Teknologi%20Pendidikan
```

## Authentication

Endpoint pencarian dan penambahan reference bersifat public karena dapat dipakai sebelum registrasi.

Endpoint profile membutuhkan:

```http
Authorization: Bearer {accessToken}
```

## Mapping Field Berdasarkan Role

| Role | Label FE | Endpoint pencarian | Field profile |
|---|---|---|---|
| `Akademisi` | Bidang Keilmuan | `GET /knowledge/search` | `expertise` |
| `Praktisi` | Bidang Industri | `GET /industries/search` | `industry` |
| `Profesional` | Bidang Keahlian | `GET /expertise/search` | `expertise` |

## Alur Pencarian Reference

1. FE memanggil endpoint `GET search` saat user mengetik.
2. FE menampilkan hasil autocomplete.
3. Jika `data` kosong, FE menampilkan tombol **Tambah Baru**.
4. FE mengirim `POST` ke endpoint reference yang sesuai.
5. FE menggunakan data hasil `POST` pada form registrasi.

FE tidak boleh mengakses database secara langsung.

## Akademisi: Bidang Keilmuan

### Search

```http
GET {{url}}/knowledge/search?q=Teknologi%20Pendidikan&limit=10
```

### Create

```http
POST {{url}}/knowledge
Content-Type: application/json
```

```json
{
  "name": "Teknologi Pendidikan"
}
```

### Response

```json
{
  "data": {
    "id": "uuid-bidang",
    "name": "Teknologi Pendidikan"
  },
  "message": "Knowledge field created successfully"
}
```

Data keilmuan sekarang flat. FE tidak perlu memakai `domain_id`, `subdomain_id`, rumpun, atau sub-rumpun.

## Praktisi: Bidang Industri

### Search

```http
GET {{url}}/industries/search?q=Pertanian&limit=10
```

### Create

```http
POST {{url}}/industries
Content-Type: application/json
```

```json
{
  "name": "Pertanian, Kehutanan dan Perikanan"
}
```

### Response

```json
{
  "data": {
    "id": "uuid-industri",
    "name": "Pertanian, Kehutanan dan Perikanan"
  },
  "message": "Industry created successfully"
}
```

## Profesional: Bidang Keahlian

Data bidang keahlian menggunakan nilai spesifikasi dari Excel sebagai nilai langsung.

### Search

```http
GET {{url}}/expertise/search?q=Article&limit=10
```

### Create

```http
POST {{url}}/expertise
Content-Type: application/json
```

```json
{
  "name": "Article"
}
```

### Response

```json
{
  "data": {
    "id": "uuid-keahlian",
    "name": "Article"
  },
  "message": "Expertise created successfully"
}
```

FE tidak perlu mengirim `category` atau `specification`. Nilai spesifikasi sudah menjadi `name`.

## Universitas atau Institusi

### Search

```http
GET {{url}}/institutions/search?q=Universitas%20Indonesia&limit=10
```

### Create

```http
POST {{url}}/institutions
Content-Type: application/json
```

```json
{
  "name": "Universitas Contoh"
}
```

Gunakan hasil response sebagai nilai `affiliation` saat registrasi atau update profile.

## Validasi dan Anti-Duplikasi

Backend menormalisasi nama sebelum menyimpan:

```text
"  Universitas   Contoh  "
```

menjadi:

```text
"Universitas Contoh"
```

Validasi nama menggunakan aturan:

- Minimal 2 karakter.
- Maksimal 255 karakter.
- Huruf dan angka diperbolehkan.
- Karakter umum seperti `&`, `.`, `,`, `-`, `/`, dan tanda kurung diperbolehkan.
- Karakter HTML atau script tidak diperbolehkan.

Database menggunakan unique index dan `ON CONFLICT`, sehingga data tidak dibuat dua kali walaupun request terjadi bersamaan.

## Registrasi

Setelah reference berhasil dipilih atau dibuat, FE mengirim registrasi:

```http
POST {{url}}/auth/register
```

Contoh Akademisi:

```json
{
  "email": "dosen@example.com",
  "password": "Secret123!",
  "role": "Akademisi",
  "firstName": "Budi",
  "lastName": "Santoso",
  "affiliation": "Universitas Contoh",
  "location": "Bandung",
  "expertise": "Teknologi Pendidikan",
  "industry": "",
  "bio": "Peneliti pendidikan"
}
```

Contoh Praktisi:

```json
{
  "email": "praktisi@example.com",
  "password": "Secret123!",
  "role": "Praktisi",
  "firstName": "Sari",
  "lastName": "Wijaya",
  "affiliation": "PT Contoh Indonesia",
  "location": "Jakarta",
  "expertise": "",
  "industry": "Pertanian, Kehutanan dan Perikanan",
  "bio": "Praktisi industri"
}
```

## Update Profile

```http
PATCH {{url}}/users/me
Authorization: Bearer {accessToken}
Content-Type: application/json
```

Semua field update bersifat opsional. Kirim hanya field yang diubah.

```json
{
  "affiliation": "Universitas Contoh",
  "expertise": "Teknologi Pendidikan",
  "industry": "",
  "bio": "Bio terbaru"
}
```

Aturan penting:

- Field tidak dikirim: nilai lama dipertahankan.
- Field dikirim sebagai `""`: nilai lama dikosongkan.
- `industry` digunakan untuk Praktisi.
- `expertise` digunakan untuk Akademisi dan Profesional.

### Contoh Axios

```js
async function updateProfile(payload) {
  const response = await api.patch("/users/me", payload);
  return response.data.data;
}
```

## Avatar

Avatar menggunakan endpoint terpisah:

```http
POST {{url}}/users/me/avatar
```

Setelah upload berhasil, gunakan URL avatar dari response sesuai kontrak endpoint. Jangan mengirim file gambar pada `PATCH /users/me`.

## Error yang Perlu Ditangani FE

| Status | Arti | Perilaku FE |
|---:|---|---|
| `200` | Request berhasil | Perbarui state dengan response |
| `201` | Reference dibuat atau existing dikembalikan | Gunakan data response |
| `400` | Nama atau payload tidak valid | Tampilkan pesan validasi |
| `401` | Token tidak valid | Refresh token atau minta login kembali |
| `409` | Data duplikat | Gunakan data yang sudah ada |
| `500` | Error server | Tampilkan pesan umum |

## Checklist Implementasi FE

- Gunakan endpoint berdasarkan role.
- Cari data sebelum menawarkan tambah baru.
- Kirim `POST` hanya jika hasil pencarian kosong.
- Simpan `id` hasil reference untuk state FE.
- Kirim nilai teks yang dipilih pada request registrasi/profile.
- Jangan menggunakan route lama `/knowledge/fields`.
- Gunakan `/knowledge/search` dan `/knowledge`.
- Jangan mengirim `category`, `domain_id`, atau `subdomain_id` pada data flat terbaru.
