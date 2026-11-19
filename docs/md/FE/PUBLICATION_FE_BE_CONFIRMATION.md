# Publication Feed — Konfirmasi Integrasi FE–BE

Dokumen ini mengonfirmasi kontrak endpoint publication yang tersedia di backend Nalakarsa.

## Base URL

Development:

```text
http://localhost:8081
```

Production Railway:

```text
https://nalakarsa-be-production.up.railway.app
```

Semua endpoint menggunakan prefix:

```text
/api/v1
```

## Authentication

Semua endpoint publication membutuhkan JWT:

```http
Authorization: Bearer <access_token>
```

Tanpa token, backend mengembalikan HTTP `401 Unauthorized`.

## 1. Get Publication Feed

```http
GET /api/v1/publications/feed
```

Contoh:

```http
GET /api/v1/publications/feed?scope=for_you&limit=12
Authorization: Bearer <access_token>
```

Query parameter:

| Parameter | Wajib | Default | Keterangan |
|---|---:|---:|---|
| `scope` | Tidak | `for_you` | Nilai: `for_you` atau `following` |
| `category` | Tidak | kosong | Filter berdasarkan kategori |
| `search` | Tidak | kosong | Mencari judul, ringkasan, atau konten |
| `cursor` | Tidak | kosong | Cursor dari response sebelumnya |
| `limit` | Tidak | `12` | Jumlah item, maksimal `50` |

Response `200 OK`:

```json
{
  "data": [
    {
      "id": "2a6d1c7e-9df6-4a9b-a7a6-3a7e8ef1a123",
      "publication_type": "CASE_STUDY",
      "title": "Implementasi IoT untuk Monitoring Tanah Pertanian",
      "summary": "Ringkasan publikasi.",
      "content": "Konten publikasi.",
      "category": "Teknologi",
      "status": "PUBLISHED",
      "cover_url": "https://example.com/cover.jpg",
      "view_count": 120,
      "save_count": 15,
      "created_at": "2026-10-22T08:00:00Z",
      "updated_at": "2026-10-22T08:00:00Z",
      "is_saved": false,
      "collaborators": [
        {
          "id": "c1a2b3c4-d5e6-7890-abcd-ef1234567890",
          "name": "Dr. Budi Santoso",
          "role": "Penulis",
          "affiliation": "Universitas Indonesia",
          "avatar_url": "https://example.com/avatar.jpg"
        }
      ]
    }
  ],
  "meta": {
    "next_cursor": "eyJjcmVhdGVkX2F0Ijoi...",
    "has_more": true,
    "limit": 12
  }
}
```

Jika sudah tidak ada data berikutnya:

```json
{
  "next_cursor": null,
  "has_more": false,
  "limit": 12
}
```

Gunakan `next_cursor` untuk request lazy scroll berikutnya. Jangan membuat cursor sendiri dari FE.

## 2. Get Publication Detail

```http
GET /api/v1/publications/{publicationId}
Authorization: Bearer <access_token>
```

Response:

```json
{
  "data": {
    "id": "2a6d1c7e-9df6-4a9b-a7a6-3a7e8ef1a123",
    "publication_type": "CASE_STUDY",
    "title": "Implementasi IoT untuk Monitoring Tanah Pertanian",
    "summary": "Ringkasan publikasi.",
    "content": "Konten lengkap publikasi.",
    "category": "Teknologi",
    "status": "PUBLISHED",
    "cover_url": "",
    "view_count": 120,
    "save_count": 15,
    "created_at": "2026-10-22T08:00:00Z",
    "updated_at": "2026-10-22T08:00:00Z",
    "is_saved": false,
    "collaborators": []
  }
}
```

Backend hanya menampilkan publication dengan status `PUBLISHED`.

## 3. Save Publication

```http
POST /api/v1/publications/{publicationId}/save
Authorization: Bearer <access_token>
```

Response:

```json
{
  "data": {
    "publication_id": "2a6d1c7e-9df6-4a9b-a7a6-3a7e8ef1a123",
    "is_saved": true,
    "save_count": 16
  }
}
```

Request save bersifat idempotent. Memanggil endpoint berkali-kali tidak menambah jumlah save berulang kali.

## 4. Unsave Publication

```http
DELETE /api/v1/publications/{publicationId}/save
Authorization: Bearer <access_token>
```

Response:

```json
{
  "data": {
    "publication_id": "2a6d1c7e-9df6-4a9b-a7a6-3a7e8ef1a123",
    "is_saved": false,
    "save_count": 15
  }
}
```

## Publication Type

Nilai yang didukung:

- `CASE_STUDY`
- `PRACTICAL_ARTICLE`
- `APPLIED_RESEARCH`
- `POLICY_BRIEF`
- `IMPLEMENTATION_REPORT`
- `IMPLEMENTATION_GUIDE`

## FE Integration Notes

1. Simpan `next_cursor` dari response feed.
2. Request halaman berikutnya menggunakan cursor tersebut.
3. Gunakan `is_saved` untuk menentukan state tombol save.
4. Setelah save/unsave berhasil, update `is_saved` dan `save_count` dari response.
5. Jangan menampilkan publication dengan status selain `PUBLISHED`.
6. `avatar_url` dapat bernilai `null`.
7. Request ke production wajib menggunakan HTTPS dan bearer token.

## Backend Status

Sudah tersedia:

- Feed `for_you`
- Feed `following`
- Filter kategori
- Pencarian publication
- Cursor pagination
- Publication detail
- Save dan unsave
- Collaborators
- AutoMigrate tabel publication
- Profile view count
- Format waktu pesan bergaya WhatsApp
- Status read/unread pesan
- Unread count pada daftar conversation
- Notifikasi pesan direct dan group
- Dokumentasi OpenAPI terbaru

## Update Terbaru: Profile View Count

Endpoint berikut sekarang menambah hitungan kunjungan profile:

```http
GET /api/v1/users/{userId}
```

Response profile memiliki `view_count` pada data profile dan `stats`:

```json
{
  "data": {
    "id": "user-uuid",
    "fullName": "Dr. Budi Santoso",
    "view_count": 42,
    "stats": {
      "connection_count": 10,
      "project_count": 3,
      "discussion_count": 5,
      "view_count": 42
    }
  }
}
```

`GET /api/v1/users/me` tidak menambah hitungan kunjungan karena itu adalah profile milik sendiri.

## Update Terbaru: Conversation dan Message

### Daftar conversation

```http
GET /api/v1/conversations
Authorization: Bearer <access_token>
```

Setiap conversation sekarang memiliki:

```json
{
  "id": "conversation-uuid",
  "last_message": "Halo, apa kabar?",
  "last_message_time": "2026-10-22T10:30:00Z",
  "unread_count": 2
}
```

### Response message

Endpoint:

```http
GET /api/v1/conversations/{conversationId}/messages
GET /api/v1/chat/contacts/{conversationId}/messages
Authorization: Bearer <access_token>
```

Item message memiliki format waktu yang lebih ringkas:

```json
{
  "id": "message-uuid",
  "sender": "me",
  "text": "Halo",
  "time": "2026-10-22T10:30:00Z",
  "time_label": "10:30",
  "status": "read",
  "is_read": true
}
```

Kemungkinan nilai `time_label`:

- Hari ini: `10:30`
- Hari sebelumnya: `Kemarin 10:30`
- Tahun yang sama: `21 Okt 10:30`
- Tahun berbeda: `21 Okt 2025 10:30`

Nilai `status` saat ini adalah `sent` atau `read`.

### Tandai pesan sudah dibaca

Panggil endpoint ini ketika user membuka conversation:

```http
PATCH /api/v1/conversations/{conversationId}/read
PATCH /api/v1/chat/contacts/{conversationId}/read
PATCH /api/v1/chat/{conversationId}/read
Authorization: Bearer <access_token>
```

Response:

```json
{
  "data": null,
  "message": "Messages marked as read"
}
```

## Update Terbaru: Notifikasi Pesan Baru

Pesan direct dan group yang baru dikirim otomatis membuat notification untuk penerima.

Ambil daftar notification:

```http
GET /api/v1/notifications
Authorization: Bearer <access_token>
```

Contoh notification pesan:

```json
{
  "id": "notification-uuid",
  "type": "message",
  "actor_id": "sender-uuid",
  "actor_name": "Dr. Budi Santoso",
  "resource_type": "conversation",
  "resource_id": "conversation-uuid",
  "title": "Pesan baru",
  "message": "Dr. Budi Santoso mengirim pesan baru",
  "read_at": null,
  "is_read": false,
  "created_at": "2026-10-22T10:30:00Z"
}
```

Untuk menandai notification sebagai sudah dibaca:

```http
PATCH /api/v1/notifications/{notificationId}/read
PATCH /api/v1/notifications/read-all
Authorization: Bearer <access_token>
```

Jika Firebase Realtime Database aktif, event realtime notification juga dikirim ke path:

```text
notifications/{userId}
```

## Checklist untuk FE

1. Tampilkan `view_count` dari response profile.
2. Tampilkan `unread_count` pada badge daftar conversation.
3. Gunakan `time_label` untuk tampilan waktu message.
4. Gunakan `status` atau `is_read` untuk indikator centang message.
5. Panggil endpoint mark-read ketika conversation dibuka.
6. Polling atau subscribe notification untuk pesan baru.
7. Gunakan `type === "message"` untuk membedakan notification pesan.
8. Gunakan `read_at` atau `is_read` untuk status notification.

Belum termasuk dalam kontrak utama:

- Endpoint membuat publication dari collaboration. Fitur ini bersifat opsional berdasarkan dokumen kebutuhan FE.
