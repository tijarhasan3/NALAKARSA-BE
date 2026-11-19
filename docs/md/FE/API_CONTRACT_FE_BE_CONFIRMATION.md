# Konfirmasi Penyesuaian BE terhadap API Contract Nalakarsa

Tanggal: 30 Oktober 2026  
Sumber acuan: `API_CONTRACT_NALAKARSA.md`  
Dokumen acuan tersebut diperlakukan sebagai kontrak kebutuhan FE-BE. Hasil di bawah membandingkan kontrak itu dengan implementasi BE saat ini; bukan semua isi dokumen dianggap sudah tersedia.

## Update implementasi BE - 30 Oktober 2026

Perubahan P0 berikut sudah diterapkan:

- `request_id` dibuat oleh middleware dan dikembalikan pada response sukses/error serta header `X-Request-ID`.
- Error autentikasi standar menggunakan code `UNAUTHENTICATED`.
- Validasi field menggunakan code `VALIDATION_ERROR` dengan HTTP `422`.
- Register/profile/auth response sudah membawa `mission` dan `email_verified`.
- Login melalui HTTP handler mengembalikan `403 EMAIL_NOT_VERIFIED` untuk user yang belum memverifikasi email.
- `GET /users/:id/stats` sudah tersedia.
- Pesan personal dan grup menerima `client_message_id` dan melakukan deduplikasi berdasarkan sender + conversation/group.
- Alias `POST /conversations` dan endpoint `/companies/search`, `/companies` sudah tersedia.

Perubahan tersebut sudah lolos `go test ./...`, `go build ./...`, dan `git diff --check`.

## Kesimpulan utama

BE sudah dapat dipakai untuk integrasi awal auth, profile, diskusi dasar, follow, request kolaborasi dasar, chat, notifikasi, upload asset Supabase, dan Firebase realtime write. Namun BE belum 100% sesuai kontrak canonical FE.

Penyesuaian paling penting sebelum FE mengunci integrasi:

1. Samakan format field canonical ke `snake_case` atau tetapkan adapter kompatibilitas secara resmi.
2. Tambahkan `request_id` pada response sukses dan error.
3. Selaraskan code `401` menjadi `UNAUTHENTICATED` atau sepakati mapping FE.
4. Lengkapi onboarding/profile: `mission`, interests, institution/company/location terstruktur, dan `email_verified`.
5. Lengkapi visibility dan status diskusi/project.
6. Tambahkan idempotensi chat menggunakan `client_message_id`.
7. Lengkapi lifecycle project, member, workspace, dan event realtime canonical.

## Status berdasarkan area

| Area | Status | Keterangan |
|---|---|---|
| Base URL `/api/v1` | READY | Router utama sudah menggunakan prefix ini. |
| Envelope sukses | PARTIAL | `data`, `message`, `timestamp`, dan pagination tersedia; `meta.request_id` belum ada. |
| Envelope error | PARTIAL | `error.code`, `message`, `details`, dan `timestamp` tersedia; `request_id` belum ada. |
| Auth register/login/refresh/logout | PARTIAL | Endpoint tersedia, tetapi payload/response belum sepenuhnya canonical. |
| Email verification | READY/PENDING ENV | Endpoint tersedia; pengiriman email tergantung Brevo/Railway. |
| Profile dan avatar | PARTIAL | Profile dan avatar tersedia; struktur institution/company/location belum sesuai kontrak. |
| Topics/interests | PARTIAL | Topics dan minimal tiga interests sudah tersedia; belum terhubung penuh ke register/onboarding final. |
| User search/follow | PARTIAL | Endpoint tersedia; pengecualian user aktif dan response canonical perlu disesuaikan. |
| Reference data | PARTIAL | Expertise, industry, knowledge, location, institution, organization tersedia dengan beberapa perbedaan path/payload. |
| Diskusi | PARTIAL | CRUD, reply, vote, field bisnis, dan media URL tersedia; visibility, tags array, output alias, dan filter kontrak belum lengkap. |
| Kolaborasi/project | PARTIAL | Project dan collaboration request dasar tersedia; lifecycle admin/member/output belum lengkap. |
| Chat personal | PARTIAL | REST + Firebase write + Supabase attachment tersedia; idempotensi client message dan response canonical belum lengkap. |
| Chat grup | PARTIAL | List/send/read dasar tersedia; create/update/member management belum lengkap. |
| Notifikasi | PARTIAL | List/read/unread/delete tersedia; action URL, actor object, dan filter kontrak belum lengkap. |
| Upload asset | READY | `POST /api/v1/assets` menggunakan Supabase Storage dan signed URL. |
| Realtime event | PARTIAL | Firebase message write sudah memakai ID deterministik, tetapi event/channel auth belum diformalkan. |

## Mismatch kontrak yang perlu disesuaikan BE

### 1. Response dan error global

Kontrak FE meminta:

```json
{
  "data": {},
  "meta": {
    "request_id": "req_..."
  }
}
```

dan error:

```json
{
  "error": {
    "code": "UNAUTHENTICATED",
    "message": "...",
    "details": {},
    "timestamp": "2026-10-30T10:00:00Z",
    "request_id": "req_..."
  }
}
```

Status BE saat ini:

- envelope `data` tersedia melalui helper response;
- pagination masih memakai `current_page`, `total_pages`, `total_items`, `limit`;
- `request_id` belum dibuat/dikembalikan;
- HTTP 401 saat ini umumnya menggunakan code `UNAUTHORIZED`, sedangkan kontrak meminta `UNAUTHENTICATED`;
- HTTP 422 belum menjadi standar handler validasi; sebagian validasi masih 400.

Penyesuaian P0: tambahkan middleware request ID, masukkan request ID ke semua response, dan tetapkan mapping error final.

### 2. Auth dan onboarding

Endpoint saat ini tersedia:

- `POST /api/v1/auth/register`
- `POST /api/v1/auth/login`
- `POST /api/v1/auth/refresh`
- `POST /api/v1/auth/logout`
- `POST /api/v1/auth/verify-email`
- `POST /api/v1/auth/resend-verification-email`
- `POST /api/v1/auth/forgot-password`
- `POST /api/v1/auth/reset-password`

Perbedaan yang perlu diperbaiki:

- Request register BE masih memakai beberapa field camelCase seperti `firstName`, `fullName`, dan `affiliation`, sementara canonical contract memakai `snake_case`.
- `mission` belum diterima pada register/update profile meskipun model user sudah memiliki field mission.
- Register masih mewajibkan `security_question` dan `security_answer`, sedangkan kontrak FE menjadikannya opsional sesuai mekanisme BE.
- Minimal tiga `topic_ids` belum divalidasi sebagai bagian dari register/onboarding final.
- Response auth perlu menambahkan `token_type`, `email_verified`, dan `verification_required`.
- Login user yang belum terverifikasi perlu memiliki code stabil `EMAIL_NOT_VERIFIED` dan detail email.

Keputusan yang disarankan: canonical request dan response memakai `snake_case`; alias camelCase dipertahankan sementara selama migrasi FE.

### 3. Profile, user search, dan statistik

Endpoint yang tersedia:

- `GET/PATCH /api/v1/users/me`
- `GET/PATCH /api/v1/auth/me`
- `GET/PATCH /api/v1/profile`
- `GET /api/v1/users/:id`
- `GET /api/v1/users/me/stats`
- `GET /api/v1/users/search`
- `GET /api/v1/users/suggestions`
- `POST /api/v1/users/:id/follow`

Gap:

- Kontrak meminta `GET /users/:user_id/stats` untuk profile user target; BE baru menyediakan stats milik sendiri.
- User response belum memiliki struktur `institution`, `company`, dan `location` object seperti kontrak.
- `mission` belum dipetakan konsisten ke response profile.
- `email_verified` belum menjadi field canonical response.
- Search perlu memastikan user yang sedang login tidak ikut dikembalikan.
- Follow saat ini toggle; FE harus menggunakan response state terbaru. Jika toggle dipertahankan, code conflict/idempotensi perlu didokumentasikan.
- Kontrak menyebut profile target privat, sedangkan route public profile BE saat ini dapat diakses tanpa token. Ini harus diputuskan dan ditegakkan di middleware.

### 4. Reference data

BE sudah memiliki beberapa endpoint referensi, tetapi FE perlu menggunakan mapping berikut:

| Kontrak FE | Kondisi BE |
|---|---|
| `GET /expertise/search` | Tersedia. |
| `POST /expertise` | Tersedia. |
| `GET /industries/search` | Perlu pastikan alias plural `/industries/search`. |
| `POST /industries` | Perlu pastikan alias plural `/industries`. |
| `GET /knowledge/search` | Tersedia. |
| `POST /knowledge` | Tersedia. |
| `GET /locations/provinces` | Tersedia melalui route location yang ada. |
| `GET /locations/cities` | Perlu validasi response `province_id` pada setiap item kota. |
| `GET /locations/search` | Perlu pastikan query `type=province|city`. |
| `GET /institutions/search` | Tersedia/ perlu verifikasi path canonical. |
| `POST /institutions` | Tersedia/ perlu verifikasi permission dan response. |
| `GET /companies/search` | Belum ada endpoint canonical; gunakan `/organizations/search` sementara atau tambahkan alias companies. |
| `POST /companies` | Belum ada endpoint canonical; gunakan `/organizations` sementara atau tambahkan alias companies. |

### 5. Diskusi

Endpoint utama sudah tersedia untuk list, detail, create, update, delete, replies, vote, dan unvote. Alias FE juga sudah ditambahkan untuk:

- `POST/DELETE /discussions/:id/upvote`
- `GET /discussions/:id/collaboration`
- `POST /discussions/:discussionId/collaboration-requests`

Gap terhadap kontrak:

- Filter `visibility`, `mine`, `topic_id`, dan `tab` belum lengkap.
- Sort `popular` dan `most_replied` belum lengkap; saat ini terutama newest/oldest.
- `tags` masih berupa string, sedangkan contract FE mencontohkan array.
- Contract menggunakan `output`, sedangkan BE menggunakan `expected_output`; perlu response alias atau canonical final.
- Field `visibility` belum menjadi field model diskusi.
- Status BE (`open`, `resolved`, `closed`) belum sama dengan status FE (`active`, `converted_to_collaboration`, `archived`).
- `collaboration_id` belum dikembalikan pada response diskusi.
- Create/update belum menegakkan seluruh field bisnis wajib dari kontrak.
- GET collaboration saat ini mengembalikan detail diskusi, bukan resource collaboration khusus.
- Delete reply perlu memastikan owner reply, owner diskusi, dan moderator sesuai permission contract.

### 6. Project dan collaboration

Endpoint dasar tersedia untuk project CRUD, members, applications, milestones, collaboration request, approval, rejection, dan withdraw.

Gap utama:

- `GET /projects` saat ini masih public, sedangkan kontrak meminta project relevan berdasarkan user login.
- Query `scope`, `visibility`, dan pagination contract belum lengkap.
- `POST /projects` saat ini belum menerima seluruh payload project yang canonical; pembuatan dari diskusi masih memakai endpoint khusus `/collaborations/from-discussion`.
- Alias `PATCH /collaborations/:id` belum tersedia.
- Project response belum memiliki background, objective, output, cover URL, visibility, owner object, dan member count secara lengkap.
- `visibility` open/private belum ditegakkan.
- Request belum mendukung `role_note` secara canonical.
- Filter request berdasarkan `project_id`, `discussion_id`, `status`, page, dan limit belum lengkap.
- Approval belum menerima payload detail project final dari kontrak dan perlu diverifikasi sebagai transaksi atomik.
- Endpoint invitation, remove member, leave, transfer admin, complete, output, cover, dan workspace belum tersedia.
- Pembuatan project dari diskusi harus dijaga idempotent berdasarkan `source_discussion_id`.

### 7. Chat dan penyimpanan

Yang sudah berjalan:

- REST conversation personal dan group.
- Read/unread personal.
- Supabase attachment untuk chat.
- Firebase Realtime write untuk pesan direct/group.
- Path Firebase sudah deterministic:
  - `chats/direct/{conversationId}/messages/{messageId}`
  - `chats/groups/{groupChatId}/messages/{messageId}`

Gap kontrak:

- `client_message_id` belum disimpan dan belum digunakan untuk deduplikasi retry.
- Response message belum canonical penuh: `conversation_id`, `sender_id`, `recipient_id`, `attachments`, dan `read_at` perlu dilengkapi.
- Query `before` dan `after` belum didukung konsisten; saat ini direct chat memakai cursor, group memakai page.
- `DELETE /conversations/:id` belum tersedia.
- Create/update group, add/remove member, update role, dan leave group belum tersedia.
- Event Firebase belum dibungkus dalam format `message.created`, `message.updated`, `message.deleted` dengan `channel` dan `data`.
- Authorization channel Firebase/RLS rules harus dikonfigurasi di Firebase; BE tidak boleh mempercayai channel dari client.
- PostgreSQL tetap menjadi source of truth REST; Firebase digunakan untuk realtime dan sinkronisasi, sesuai rekomendasi kontrak.

### 8. Notifikasi

Endpoint tersedia:

- `GET /api/v1/notifications`
- `GET /api/v1/notifications/unread-count`
- `POST/PATCH /api/v1/notifications/:id/read`
- `PATCH /api/v1/notifications/read-all`
- `DELETE /api/v1/notifications/:id`
- `DELETE /api/v1/notifications`

Gap:

- Filter `unread_only`, page, dan limit perlu diverifikasi.
- Response saat ini memakai `actor_id`, `actor_name`, `resource_type`, dan `resource_id`; contract FE menginginkan object `actor` serta `action_url`.
- Event notifikasi belum menjamin seluruh type contract: follow, collaboration request, approval/rejection, reply, message, dan group member added.
- Routing notification harus memiliki `action_url` atau mapping resource yang konsisten.

### 9. Upload asset dan keamanan

BE sekarang menyediakan:

```http
POST /api/v1/assets
```

Asset disimpan pada Supabase dengan pola:

```text
users/{userID}/{category}/{uuid}.{extension}
```

Attachment chat dipisahkan ke bucket `chat-attachments`, sedangkan avatar tetap di bucket avatar.

Yang masih perlu dipastikan sebelum production:

- Bucket Supabase `assets`, `avatars`, dan `chat-attachments` sudah dibuat.
- Bucket private menggunakan signed URL.
- Limit MIME dan ukuran sudah sama antara `.env`, Railway, dan FE.
- Upload rate limit dan malware scanning ditentukan.
- Endpoint upload media diskusi memakai `/assets` sebelum menyimpan `media_url`.
- File executable dan MIME spoofing ditolak; validasi MIME sebaiknya dilengkapi pemeriksaan signature file.

## Prioritas implementasi BE

### P0 - Wajib agar FE dapat mengunci kontrak

1. Middleware `request_id` dan penambahan request ID ke seluruh response.
2. Finalisasi error code `UNAUTHENTICATED`, `VALIDATION_ERROR`, dan `EMAIL_NOT_VERIFIED`.
3. Tambahkan `mission`, `email_verified`, `verification_required`, dan mapping response canonical.
4. Dukung snake_case canonical dengan alias camelCase sementara.
5. Tambahkan diskusi `visibility`, `output`, tags array, status final, dan `collaboration_id`.
6. Tambahkan `client_message_id` untuk direct/group message.
7. Sediakan `GET /users/:id/stats` dan filter request collaboration.
8. Pastikan endpoint `/companies` dan `/industries` sesuai path yang FE gunakan.

### P1 - Fitur kolaborasi dan chat lengkap

1. Project visibility open/private dan permission guest/user/member/admin.
2. Approval atomik, invitations, remove member, leave, transfer admin.
3. Complete project, output, cover, read-only workspace.
4. Group create/update/member management.
5. Canonical realtime event dan Firebase channel authorization.
6. Notification actor object, action URL, dan seluruh event bisnis.

### P2 - Hardening production

1. Audit log perubahan permission, project, member, status, dan output.
2. Idempotency key untuk approval, project creation, upload, dan vote.
3. Malware scanning dan file signature validation.
4. Observability request ID, structured logs, metrics, dan tracing.
5. Contract test otomatis dari OpenAPI untuk FE-BE.

## Kontrak yang dapat dipakai FE sekarang

FE dapat mulai memakai endpoint berikut tanpa menunggu batch P0:

- Auth dasar dan verification email.
- `/users/me`, `/auth/me`, `/profile`.
- `/profile/interests` dan `/topics`.
- `/users`, `/users/search`, followers/following, follow toggle.
- Diskusi CRUD, replies, `/votes`, `/upvote`.
- `/assets` untuk upload file/gambar.
- `/conversations` dan `/chats/groups` untuk chat dasar.
- Notification list, unread count, read, read-all, delete.
- Publication feed, detail, save, unsave.

FE perlu menunggu konfirmasi kontrak P0 untuk menggunakan secara penuh:

- onboarding final;
- visibility private;
- project/member lifecycle;
- workspace file/read-only;
- deduplikasi chat;
- canonical error/request ID.

## Konfirmasi yang diperlukan dari FE

1. Apakah FE akan menjadikan `snake_case` sebagai format utama dan mempertahankan camelCase hanya sebagai fallback?
2. Apakah profile user lain wajib login, atau boleh dibaca guest untuk profile publik?
3. Apakah Firebase tetap menjadi provider realtime production, dengan PostgreSQL sebagai source of truth REST?
4. Apakah FE membutuhkan `after` cursor untuk chat atau cukup `before`/pagination page?
5. Apakah tags dikirim sebagai array string pada kontrak final?
6. Apakah `output` akan menggantikan `expected_output`, atau keduanya dipertahankan sebagai alias?
7. Apakah upload asset diskusi selalu melalui `POST /assets` sebelum create/update diskusi?

## Status verifikasi

Perubahan BE sebelumnya sudah diverifikasi dengan:

- `go test ./...`
- `go build ./...`
- `git diff --check`

Dokumen ini adalah hasil audit kontrak terhadap kode saat ini. Belum ada perubahan kode tambahan yang dilakukan pada audit ini.

