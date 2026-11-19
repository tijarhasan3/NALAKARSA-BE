# NALAKARSA — Konfirmasi Kebutuhan FE–BE

Dokumen ini adalah hasil pencocokan kebutuhan pada `NALAKARSA_BE_COMMUNICATION.md` dengan backend saat ini.

## Status Implementasi Terakhir — 30 Oktober 2026

Penyesuaian BE tahap utama sudah diterapkan dan sudah dimigrasikan ke database lokal `nalakarsa`.

Validasi lokal:

```text
PostgreSQL: localhost:5432
Database: nalakarsa
AutoMigrate: berhasil
go test ./...: berhasil
go build ./...: berhasil
```

Kolom baru pada tabel `discussions`:

```text
background
objective
expected_output
media_url
topic_id
```

Status:

- **READY**: endpoint dan kemampuan utama sudah tersedia.
- **PARTIAL**: tersedia, tetapi path, payload, permission, enum, atau response belum sama dengan kebutuhan FE.
- **GAP**: belum tersedia dan perlu perubahan BE.

## 1. Kontrak Umum

### Base URL

```text
Development: http://localhost:8081/api/v1
Production:  https://nalakarsa-be-production.up.railway.app/api/v1
```

### Authentication

Endpoint protected menggunakan:

```http
Authorization: Bearer <access_token>
```

### Response sukses saat ini

Sebagian besar endpoint menggunakan wrapper:

```json
{
  "data": {},
  "meta": {},
  "message": "Success",
  "timestamp": "2026-10-30T00:00:00Z"
}
```

Publication feed adalah pengecualian: response-nya langsung memiliki `data` dan `meta` sesuai kontrak feed.

### Response error saat ini

```json
{
  "error": {
    "code": "VALIDATION_ERROR",
    "message": "Validation failed",
    "details": {},
    "timestamp": "2026-10-30T00:00:00Z"
  }
}
```

**Perlu keputusan:** default error BE saat ini memakai `INTERNAL_ERROR`, sedangkan dokumen FE menyebut `SERVER_ERROR`. FE sementara perlu menerima keduanya atau BE menyamakan menjadi satu kode final.

## 2. Auth dan Onboarding

| Kebutuhan FE | Endpoint BE saat ini | Status | Catatan |
|---|---|---|---|
| Register | `POST /auth/register` | PARTIAL | Profil profesional tersedia, tetapi belum menerima minimal tiga topik minat. |
| Verify email | `POST /auth/verify-email` | READY | Kode 6 digit dikendalikan BE. |
| Resend verification | `POST /auth/resend-verification-email` | READY | Nama path berbeda dari usulan FE. |
| Login | `POST /auth/login` | READY | Mengembalikan access token, refresh token, dan user. |
| Refresh | `POST /auth/refresh` | READY |  |
| Logout | `POST /auth/logout` | READY | Protected endpoint. |
| Auth me | `GET /auth/me` | READY | Alias dari `GET /users/me`. |
| Profile | `GET /profile` | READY | Alias baru dari `GET /users/me`. |
| Update profile | `PATCH /profile` | READY | Alias baru dari `PATCH /users/me`. |
| Profile interests | `PATCH /profile/interests` | GAP | Belum ada endpoint onboarding/interests khusus. |
| Topics | `GET /topics` | GAP | Belum ada resource topics terpisah. |
| Topic suggestions | `POST /topics/suggestions` | GAP | Belum tersedia. |

### Penyesuaian yang dibutuhkan

1. Tambahkan `onboarding_status` atau `onboarding_completed` pada response auth/me.
2. Tentukan apakah minat menggunakan resource `skills`, resource `topics`, atau kombinasi keduanya.
3. Topics dan interests masih perlu resource/endpoint final.

## 3. Diskusi

### Endpoint

| Kebutuhan FE | Endpoint BE saat ini | Status |
|---|---|---|
| List | `GET /discussions` | PARTIAL |
| Detail | `GET /discussions/{id}` | READY |
| Create | `POST /discussions` | PARTIAL |
| Update | `PATCH /discussions/{id}` | PARTIAL |
| Delete | `DELETE /discussions/{id}` | READY |
| Upvote | `POST /discussions/{id}/votes` | READY |
| Unvote | `DELETE /discussions/{id}/votes` | READY |
| Replies list | `GET /discussions/{id}/replies` | READY |
| Add reply | `POST /discussions/{id}/replies` | READY |
| Delete reply | `DELETE /discussions/{id}/replies/{reply_id}` | READY |
| Collaboration marker | `POST /discussions/{id}/collaboration` | PARTIAL |
| Collaboration detail | `GET /discussions/{id}/collaboration` | GAP |

### Field mismatch

Response diskusi saat ini memiliki:

```text
title, description, category, tags, status,
isInCollaboration, replies, upvote_count, has_upvoted,
author, role, avatar_url, createdAt
```

Dokumen FE membutuhkan tambahan atau mapping final untuk:

```text
topic_id
background
objective
expected_output
media_url
has_active_collaboration
collaboration_id
updated_at
```

BE perlu memilih salah satu:

1. Menambah field tersebut pada model/response `Discussion`; atau
2. Menetapkan mapping resmi, misalnya `description` menjadi `background`.

### Filter dan pagination

BE saat ini mendukung `search/q`, `category`, `role`, `sort`, `page`, dan `limit`.

Belum tersedia:

- `tab=for_you|all`.
- `topic` atau `topic_id`.
- Cursor pagination.
- Ranking khusus `for_you`.

Untuk fase awal FE dapat menggunakan page-based pagination. Jika lazy scroll berbasis cursor wajib, BE perlu menambah cursor contract.

## 4. Permintaan Kolaborasi

### Endpoint saat ini

```http
POST   /collaborations/requests
GET    /collaborations/requests
PATCH  /collaborations/requests/{id}/approve
PATCH  /collaborations/requests/{id}/reject
DELETE /collaborations/requests/{id}
```

Status: **PARTIAL**.

Yang sudah tersedia:

- Pengajuan berdasarkan `discussionId` atau `projectId`.
- `proposedContribution`.
- Status pending, accepted, rejected.
- Approve/reject oleh pemilik diskusi atau project.
- Notifikasi request dan hasil approve/reject.
- Approve membuat project dan group chat.

Perbedaan dari kebutuhan FE:

- Path menggunakan `/collaborations/requests`, bukan `/discussions/{id}/collaboration-requests`.
- Payload belum memiliki `role_id`, `role_note`, dan `contribution` terpisah.
- Belum ada `GET /collaboration-requests/{requestId}`.
- Belum ada aturan eksplisit bahwa request rejected tidak boleh diajukan kembali.
- Belum ada flow join request yang berbeda dari collaboration request.

## 5. Kolaborasi / Project

BE saat ini menggunakan resource `projects` sebagai representasi kolaborasi:

```http
GET    /projects
POST   /projects
GET    /projects/{id}
PATCH  /projects/{id}
DELETE /projects/{id}
GET    /projects/{id}/members
```

Status: **PARTIAL**.

Yang sudah ada:

- List dan detail project.
- Create/update/delete project.
- Members, applications, dan milestones.
- Status project: `draft`, `open`, `in_review`, `active`, `completed`, `archived`.
- Collaboration request dapat membuat project dan group chat.
- Project dapat dibuat dari discussion melalui `POST /collaborations/from-discussion`.

Yang belum tersedia:

- `cover_url` dan upload cover.
- `visibility=open|private`.
- Output collaboration: description, URL, file, visibility.
- Invitations.
- Join request khusus collaboration.
- Remove member.
- Leave collaboration.
- Transfer admin.
- Complete collaboration khusus.
- Read-only enforcement setelah completed.
- `completed_at`.

### Gap penting

Flow membuat project dari discussion sekarang tersedia melalui:

```http
POST /collaborations/from-discussion
```

Payload:

```json
{ "discussionId": "discussion-uuid" }
```

## 6. Workspace, Chat, dan File

### Chat yang tersedia

Direct chat:

```http
GET   /conversations
POST  /conversations/direct
GET   /conversations/{id}/messages
POST  /conversations/{id}/messages
PATCH /conversations/{id}/read
```

Group chat:

```http
GET  /chats/groups
GET  /chats/groups/{groupId}/messages
POST /chats/groups/{groupId}/messages
```

Alias mark-read:

```http
PATCH /chat/contacts/{id}/read
PATCH /chat/{id}/read
```

Status: **PARTIAL** untuk workspace collaboration.

Yang sudah tersedia:

- Direct chat dan group chat.
- Upload attachment chat.
- Delete message milik sendiri.
- Validasi membership.
- `unread_count` pada daftar direct conversation.
- Message memiliki `time_label`, `status`, dan `is_read`.
- Notifikasi pesan baru direct/group.

Yang belum tersedia:

- Workspace endpoint berdasarkan `collaborationId`.
- Daftar file workspace.
- Delete file workspace.
- Admin menghapus pesan/file siapa pun.
- Read-only behavior setelah collaboration completed.
- Read/unread khusus group chat.
- Unread count terpisah untuk direct dan workspace group.

### Format message terbaru

```json
{
  "id": "message-uuid",
  "sender": "me",
  "text": "Halo",
  "time": "2026-10-30T10:30:00Z",
  "time_label": "10:30",
  "status": "read",
  "is_read": true
}
```

## 7. Beranda dan Carousel

Endpoint yang tersedia:

```http
GET /dashboard
GET /dashboard/summary
GET /dashboard/ongoing-projects
GET /dashboard/notifications
GET /collaborations
GET /projects
```

Status: **PARTIAL**.

Pagination saat ini menggunakan:

```json
{
  "current_page": 1,
  "total_pages": 2,
  "total_items": 20,
  "limit": 10
}
```

Dokumen FE mengusulkan:

```json
{
  "next_cursor": null,
  "has_next": false
}
```

Rekomendasi: gunakan page-based untuk endpoint lama dan cursor-based hanya untuk feed baru yang membutuhkan lazy scroll.

Filter `section=mine|completed|active` dan `topic_id` global belum tersedia sebagai kontrak khusus.

## 8. Notifikasi

### Endpoint saat ini

```http
GET   /notifications
PATCH /notifications/{id}/read
PATCH /notifications/read-all
GET   /notifications/unread-count
POST  /notifications/{id}/read
DELETE /notifications/{id}
DELETE /notifications
```

Status: **PARTIAL**.

Yang sudah tersedia:

- List notification dengan pagination.
- Mark one as read dan mark all as read.
- `type`, `title`, `message`, `actor_id`, `actor_name`, `resource_type`, `resource_id`, `read_at`, `is_read`, `created_at`.
- Notification follow, collaboration, dan message pada flow tertentu.

Belum tersedia:

- Object `actor` penuh; saat ini menggunakan `actor_id` dan `actor_name`.
- Event mention.
- Event perubahan admin dan perubahan penting collaboration secara menyeluruh.

Chat notification menggunakan indikator terpisah:

```text
type: message
resource_type: conversation atau group_chat
is_read: false
```

## 9. Permission Matrix — Kondisi Saat Ini

| Aksi | Kondisi BE saat ini |
|---|---|
| Melihat diskusi/project detail | Umumnya public |
| Membuat diskusi/project | User login |
| Edit/hapus diskusi | Pemilik |
| Upvote/reply | User login |
| Submit collaboration request | User login, bukan pemilik resource |
| Approve/reject request | Initiator/pemilik resource |
| Melihat project members | Saat ini endpoint public; perlu dikonfirmasi |
| Chat direct | Member conversation |
| Chat group | Member group yang valid |
| Upload attachment | Member conversation/group |
| Hapus message | Pengirim sendiri; admin-delete semua belum tersedia |
| Invite/remove/leave/transfer admin | Belum tersedia |

## 10. Error Contract

Kode yang sudah digunakan:

```text
UNAUTHORIZED
FORBIDDEN
NOT_FOUND
VALIDATION_ERROR
CONFLICT
RATE_LIMITED
INTERNAL_ERROR
```

Yang perlu disesuaikan:

- Sepakati `SERVER_ERROR` versus `INTERNAL_ERROR`.
- Tambahkan `FILE_TOO_LARGE`.
- Tambahkan `UNSUPPORTED_FILE_TYPE`.
- Pastikan conflict selalu memakai HTTP `409`.
- Pastikan endpoint mutasi mengembalikan resource terbaru, bukan hanya `{id}` atau `null`.

## 11. Jawaban atas Pertanyaan FE

| Pertanyaan | Jawaban sementara |
|---|---|
| Collaboration dibuat bagaimana? | Saat ini melalui `POST /projects` atau approve collaboration request. Flow khusus dari discussion belum memiliki route terdaftar. |
| Cover default dari BE? | Belum tersedia. Cover URL/upload perlu ditambahkan. |
| Upload file resmi? | Attachment chat tersedia; file workspace belum tersedia. |
| Output akhir? | Belum ada kontrak output collaboration. |
| Enum status final? | Discussion memakai `open/resolved/closed`; project memakai beberapa status. Mapping final perlu disepakati. |
| Sorting/ranking? | Sort tersedia sebagian; ranking `for_you` belum tersedia. |
| Tab Untuk Anda? | Belum ada parameter `tab`; perlu topic/interests atau ranking baru. |
| Admin/member/pending? | Tersedia sebagian melalui members, applications, dan requests. |
| Notifikasi perubahan collaboration? | Request approve/reject tersedia; event admin/perubahan penting belum lengkap. |
| Realtime chat? | Firebase push digunakan jika konfigurasi tersedia; REST tetap sumber utama. |
| Unread count personal/workspace? | Personal tersedia; group/workspace terpisah belum tersedia. |
| Audit trail/delete policy? | Belum ada audit trail lengkap untuk file, member, dan notification. |

## 12. Prioritas Penyesuaian BE

### P0 — Sebelum FE batch utama

### P0 — Selesai

1. Alias endpoint profile.
2. Status onboarding dasar pada response user.
3. Field tambahan diskusi dan AutoMigrate lokal.
4. Route membuat project dari discussion.
5. Unread count dan delete notification.

### P0 — Masih perlu keputusan

1. Finalisasi endpoint onboarding topics/interests.
2. Hubungkan `collaboration_id` diskusi ke project aktif.
3. Finalisasi status dan permission collaboration.
4. Samakan error code `SERVER_ERROR` atau `INTERNAL_ERROR`.

### P1 — Collaboration lengkap

1. Cover collaboration.
2. Invitation, join request, approve/reject, remove member, leave, transfer admin.
3. Output collaboration dan completed timestamp.
4. Workspace files dan permission admin.
5. Read-only behavior setelah completed.

### P2 — Penyempurnaan

1. Cursor pagination untuk diskusi dan collaboration.
2. Ranking `for_you` berbasis topics/interests.
3. Mention dan event perubahan penting.
4. Audit trail.

## 13. Endpoint yang Bisa Dipakai FE Sekarang

```text
POST   /api/v1/auth/register
POST   /api/v1/auth/login
POST   /api/v1/auth/refresh
POST   /api/v1/auth/verify-email
POST   /api/v1/auth/resend-verification-email
GET    /api/v1/auth/me
GET    /api/v1/users/me
PATCH  /api/v1/users/me
GET    /api/v1/profile
PATCH  /api/v1/profile

GET    /api/v1/discussions
GET    /api/v1/discussions/{id}
POST   /api/v1/discussions
PATCH  /api/v1/discussions/{id}
DELETE /api/v1/discussions/{id}
POST   /api/v1/discussions/{id}/votes
DELETE /api/v1/discussions/{id}/votes
GET    /api/v1/discussions/{id}/replies
POST   /api/v1/discussions/{id}/replies

GET    /api/v1/projects
GET    /api/v1/projects/{id}
GET    /api/v1/projects/{id}/members
POST   /api/v1/projects/{id}/applications
GET    /api/v1/collaborations/requests
POST   /api/v1/collaborations/requests
PATCH  /api/v1/collaborations/requests/{id}/approve
PATCH  /api/v1/collaborations/requests/{id}/reject
POST   /api/v1/collaborations/from-discussion

GET    /api/v1/conversations
GET    /api/v1/conversations/{id}/messages
POST   /api/v1/conversations/{id}/messages
PATCH  /api/v1/conversations/{id}/read
GET    /api/v1/chats/groups
GET    /api/v1/chats/groups/{groupId}/messages
POST   /api/v1/chats/groups/{groupId}/messages

GET    /api/v1/notifications
GET    /api/v1/notifications/unread-count
PATCH  /api/v1/notifications/{id}/read
POST   /api/v1/notifications/{id}/read
PATCH  /api/v1/notifications/read-all
DELETE /api/v1/notifications/{id}
DELETE /api/v1/notifications
```

## 14. Kesimpulan

Backend sudah mencakup fondasi auth, diskusi, project/collaboration request, chat, notification, publication, profile view count, serta read/unread message.

Backend belum dapat dianggap **100% ready** untuk seluruh dokumen FE karena onboarding topics, field detail diskusi, lifecycle collaboration, invitation/join flow, workspace files, pagination cursor, dan beberapa kontrak error masih perlu difinalisasi.

FE dapat mulai dari endpoint pada bagian 13. Untuk fitur berstatus `PARTIAL` atau `GAP`, FE sebaiknya menunggu keputusan kontrak final sebelum membuat UI yang bergantung pada field atau permission tersebut.

