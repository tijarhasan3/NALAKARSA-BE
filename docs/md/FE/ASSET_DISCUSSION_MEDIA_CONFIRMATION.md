# Konfirmasi FE–BE: Upload Gambar dan `media_url` Diskusi

Tanggal: 17 November 2026  
Status: BE sudah disesuaikan

## Masalah yang ditemukan

File berhasil di-upload ke Supabase, tetapi response create discussion masih mengembalikan:

```json
{
  "media_url": ""
}
```

Penyebabnya: upload asset dan create discussion adalah dua request berbeda. URL hasil upload harus dikirim kembali pada request create discussion.

## Alur yang harus digunakan FE

### 1. Upload gambar

```http
POST /api/v1/assets
Authorization: Bearer <access_token>
Content-Type: multipart/form-data
```

Field multipart:

```text
file: <gambar>
category: discussion
```

Contoh response:

```json
{
  "data": {
    "path": "users/<user_id>/discussion/<uuid>.png",
    "url": "https://<supabase-project>.supabase.co/storage/v1/object/sign/assets/users/<user_id>/discussion/<uuid>.png?...",
    "name": "gambar.png",
    "mime_type": "image/png",
    "size": 123456,
    "category": "discussion"
  },
  "message": "Asset uploaded successfully"
}
```

FE wajib mengambil nilai:

```js
const mediaUrl = uploadResponse.data.url
```

### 2. Kirim URL saat membuat diskusi

```http
POST /api/v1/discussions
Authorization: Bearer <access_token>
Content-Type: application/json
```

Payload:

```json
{
  "title": "Aktualisasi Mesin Produksi",
  "description": "Deskripsi diskusi",
  "background": "Latar belakang diskusi",
  "objective": "Berkolaborasi",
  "expected_output": "Prototype",
  "category": "Riset",
  "tags": null,
  "media_url": "https://<supabase-project>.supabase.co/storage/v1/object/sign/assets/users/<user_id>/discussion/<uuid>.png?..."
}
```

BE juga menerima alias berikut jika FE sudah terlanjur menggunakannya:

```json
{
  "asset_url": "https://..."
}
```

atau:

```json
{
  "image_url": "https://..."
}
```

Nilai tersebut akan disimpan ke kolom `media_url`.

## Response create discussion

Jika URL sudah dikirim, response akan berisi:

```json
{
  "data": {
    "id": "711531e0-dd71-48d9-855b-0fbed65da04a",
    "title": "Aktualisasi Mesin Produksi",
    "media_url": "https://<supabase-project>.supabase.co/storage/v1/object/sign/assets/users/<user_id>/discussion/<uuid>.png?...",
    "category": "Riset",
    "tags": "",
    "status": "open"
  }
}
```

## Status bucket Supabase

Bucket yang digunakan saat ini bersifat **PUBLIC**:

| Bucket | Status | Kegunaan |
|---|---|---|
| `assets` | PUBLIC | Gambar/file diskusi dan asset umum |
| `Chat-attachments` | PUBLIC | Lampiran pesan |
| `avatars` | PUBLIC | Foto profil |

ENV Railway harus menggunakan nama bucket yang sama persis:

```env
SUPABASE_ASSET_BUCKET=assets
SUPABASE_CHAT_BUCKET=Chat-attachments
SUPABASE_STORAGE_BUCKET=avatars
```

Karena bucket PUBLIC, URL object dapat diakses melalui URL public Supabase. BE saat ini tetap mengembalikan signed URL pada `data.url` agar format response tetap konsisten. FE cukup menggunakan `data.url` sebagai sumber gambar.

## `tags` bersifat opsional

FE dapat mengirim salah satu bentuk berikut:

```json
"tags": null
```

```json
"tags": []
```

```json
"tags": ["iot", "mesin produksi"]
```

atau tidak mengirim field `tags`. Semua format tersebut diterima BE.

## Checklist FE

- Tunggu response upload `/api/v1/assets` selesai.
- Ambil URL dari `response.data.url`, bukan `response.data.path`.
- Kirim URL tersebut sebagai `media_url` pada request create discussion.
- Jangan mengirim object response upload secara langsung sebagai `media_url`.
- Jangan hanya mengirim `path` jika FE membutuhkan URL yang dapat ditampilkan.
- Pastikan deployment Railway sudah berstatus `ACTIVE`.

## Catatan signed URL

URL asset bersifat signed URL dan memiliki masa berlaku. Untuk menampilkan asset setelah URL kedaluwarsa, FE perlu meminta signed URL baru menggunakan endpoint asset/detail yang tersedia atau BE perlu menyediakan endpoint refresh URL.

