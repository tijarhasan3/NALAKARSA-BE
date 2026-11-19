# Konfirmasi Penyesuaian BE terhadap Business Summary NALAKARSA

Tanggal: 30 Oktober 2026  
Sumber acuan: `NALAKARSA_BUSINESS_SUMMARY.md`  
Dokumen ini adalah konfirmasi teknis untuk FE. Dokumen bisnis diperlakukan sebagai acuan kebutuhan, bukan instruksi untuk mengubah seluruh fitur secara otomatis.

## Ringkasan status

| Area | Status BE | Keterangan untuk FE |
|---|---|---|
| Auth, login, refresh, logout | READY | Endpoint auth utama sudah tersedia. |
| Verifikasi email dan reset password | READY/PENDING ENV | Kontrak tersedia; pengiriman email bergantung pada konfigurasi Brevo di deployment. |
| Profile dasar dan statistik | READY | Gunakan `/profile` atau `/users/me`; `view_count` sudah dihitung saat profil publik dibuka. |
| Onboarding role/profil | PARTIAL | Status onboarding dasar tersedia; update topik minat sudah tersedia, tetapi validasi topik saat registrasi dan status onboarding final belum selesai. |
| Topik resmi dan usulan topik | PARTIAL | Endpoint list, search, suggestion, dan update interests tersedia; approval admin belum tersedia. |
| Diskusi publik | READY | CRUD, detail, replies, vote, dan field bisnis utama tersedia. |
| Link diskusi ke kolaborasi | PARTIAL | Diskusi dapat ditandai dan dapat dibuat menjadi kolaborasi, tetapi `collaboration_id` belum dikembalikan. |
| Permintaan kolaborasi | PARTIAL | Submit, list, approve, reject, withdraw tersedia; aturan role dan lifecycle anggota belum lengkap. |
| Katalog kolaborasi | PARTIAL | Listing project tersedia, tetapi filter global, sorting bisnis, section, dan hak akses final belum lengkap. |
| Workspace chat | PARTIAL | Group chat tersedia; aturan read-only setelah selesai dan file workspace belum lengkap. |
| Asset umum | READY | `POST /api/v1/assets` menyimpan file/gambar ke Supabase bucket `assets` dengan signed URL. |
| File workspace | GAP | Belum ada endpoint file khusus workspace. |
| Penyelesaian kolaborasi | GAP/PARTIAL | Model dan endpoint final untuk completed, output, visibility, dan read-only belum lengkap. |
| Notifikasi | PARTIAL/READY | List, unread count, mark read, delete satu, dan delete semua tersedia; semua event bisnis belum terhubung. |
| Publications | READY | Feed, detail, save, unsave tersedia; fitur ini tidak dirinci dalam Business Summary. |

## Kontrak BE yang dapat langsung dipakai FE

### Auth

| Method | Endpoint | Kegunaan |
|---|---|---|
| POST | `/api/v1/auth/register` | Registrasi user. |
| POST | `/api/v1/auth/login` | Login dan memperoleh token. |
| POST | `/api/v1/auth/refresh` | Refresh token. |
| POST | `/api/v1/auth/logout` | Logout user terautentikasi. |
| POST | `/api/v1/auth/verify-email` | Verifikasi email. |
| POST | `/api/v1/auth/resend-verification-email` | Kirim ulang email verifikasi. |
| POST | `/api/v1/auth/forgot-password` | Meminta reset password. |
| POST | `/api/v1/auth/reset-password` | Mengubah password dengan token reset. |

### Profile dan user

Gunakan salah satu alias berikut secara konsisten; rekomendasi untuk FE adalah `/profile` untuk halaman profile milik sendiri.

| Method | Endpoint | Kegunaan |
|---|---|---|
| GET | `/api/v1/profile` | Profile user yang sedang login. |
| PATCH | `/api/v1/profile` | Update profile user yang sedang login. |
| GET | `/api/v1/profile/interests` | Mengambil topik minat user. |
| PATCH | `/api/v1/profile/interests` | Mengganti topik minat; minimal tiga topic ID. |
| GET | `/api/v1/users/me/stats` | Statistik profile milik sendiri. |
| GET | `/api/v1/users/:id` | Profile publik dan hitungan kunjungan. |
| GET | `/api/v1/users/:id/followers` | Daftar followers. |
| GET | `/api/v1/users/:id/following` | Daftar following. |
| POST | `/api/v1/users/:id/follow` | Follow atau unfollow user. |

### Topics

| Method | Endpoint | Kegunaan |
|---|---|---|
| GET | `/api/v1/topics` | Daftar topik resmi, mendukung `search`. |
| POST | `/api/v1/topics/suggestions` | Mengusulkan topik baru; status awal `pending`. |

Response profile sekarang memiliki field tambahan:

```json
{
  "onboarding_status": "completed",
  "onboarding_completed": true,
  "view_count": 12
}
```

Nilai `onboarding_status` saat ini baru membedakan `email_verification_pending` dan `completed`. Status berbasis kelengkapan topik/profil belum dapat dianggap final sebelum modul topik selesai.

### Diskusi

| Method | Endpoint | Kegunaan |
|---|---|---|
| GET | `/api/v1/discussions` | Feed diskusi publik. |
| GET | `/api/v1/discussions/:id` | Detail diskusi. |
| POST | `/api/v1/discussions` | Membuat diskusi. |
| PATCH | `/api/v1/discussions/:id` | Mengubah diskusi sebelum menjadi kolaborasi aktif. |
| DELETE | `/api/v1/discussions/:id` | Menghapus diskusi sesuai permission. |
| GET | `/api/v1/discussions/:id/replies` | Daftar jawaban. |
| POST | `/api/v1/discussions/:id/replies` | Menambah jawaban. |
| POST | `/api/v1/discussions/:id/vote` | Upvote atau membatalkan vote. |
| POST | `/api/v1/discussions/:id/collaboration` | Menandai diskusi sebagai kolaborasi pada kontrak lama. |

Payload create/update diskusi kini mendukung:

```json
{
  "title": "Judul diskusi",
  "topic_id": "uuid-topik",
  "background": "Latar belakang kolaborasi",
  "objective": "Tujuan kolaborasi",
  "expected_output": "Output yang diharapkan",
  "media_url": "https://..."
}
```

Catatan: `topic_id` sudah disediakan di model dan request DTO, tetapi endpoint topik resmi belum tersedia. FE jangan menganggap `topic_id` dapat diisi dari daftar topik server sampai modul Topics selesai.

### Kolaborasi dan permintaan anggota

| Method | Endpoint | Kegunaan |
|---|---|---|
| GET | `/api/v1/collaborations` | Listing project/kolaborasi. |
| GET | `/api/v1/projects/:id` | Detail kolaborasi. |
| GET | `/api/v1/projects/:id/members` | Daftar anggota. |
| POST | `/api/v1/collaborations/requests` | Mengajukan kolaborasi/bergabung. |
| POST | `/api/v1/discussions/:discussionId/collaboration-requests` | Alias FE untuk mengajukan kolaborasi langsung dari detail diskusi. |
| GET | `/api/v1/collaborations/requests` | Melihat request yang relevan. |
| PATCH | `/api/v1/collaborations/requests/:id/approve` | Menyetujui request. |
| PATCH | `/api/v1/collaborations/requests/:id/reject` | Menolak request. |
| DELETE | `/api/v1/collaborations/requests/:id` | Menarik request. |
| POST | `/api/v1/collaboration-requests/:id/approve` | Alias FE untuk menyetujui request. |
| POST | `/api/v1/collaboration-requests/:id/reject` | Alias FE untuk menolak request. |
| POST | `/api/v1/collaborations/from-discussion` | Membuat kolaborasi dari diskusi yang dipilih. |

Contoh payload membuat kolaborasi dari diskusi:

```json
{
  "discussionId": "uuid-diskusi"
}
```

Status saat ini belum sepenuhnya sama dengan alur bisnis pada Business Summary. Endpoint tersebut belum menjadi satu transaksi lengkap yang sekaligus mengisi cover, tipe akses, admin pertama, anggota pengaju pertama, dan workspace.

## Gap BE yang harus diimplementasikan

### P0 - Wajib sebelum FE mengunci alur utama

1. **Modul Topics**

   - Model topik resmi.
   - Endpoint list/search topik untuk onboarding, diskusi, filter, dan rekomendasi.
   - Relasi user-topik minat.
   - Validasi minimal tiga topik saat onboarding.
   - Update topik minat melalui profile.
   - Endpoint usulan topik baru dan status review admin.

2. **Kontrak onboarding final**

   - Validasi role, bidang keahlian, instansi/universitas, provinsi, kota/kabupaten, misi, dan topik.
   - Status yang membedakan minimal `incomplete`, `email_verification_pending`, dan `completed`.
   - Penolakan akses ke endpoint aplikasi jika onboarding wajib belum selesai.
   - Penetapan apakah draft onboarding hanya dikelola FE atau juga disimpan BE.

3. **Relasi diskusi-kolaborasi**

   - Tambahkan relasi nyata dari satu diskusi ke satu kolaborasi aktif.
   - Kembalikan `collaboration_id`, `collaboration_status`, dan URL/detail kolaborasi pada response diskusi.
   - Kunci field utama diskusi setelah kolaborasi aktif.
   - Pengajuan berikutnya diarahkan sebagai request bergabung, bukan membuat kolaborasi kedua.

4. **Hak akses guest**

   - Konfirmasi aturan final: Business Summary menyatakan guest tidak dapat membuka profile atau katalog kolaborasi.
   - Jika aturan ini final, endpoint profile publik dan listing kolaborasi harus membutuhkan login atau dipisahkan menjadi endpoint publik yang disetujui.
   - FE tidak boleh hanya mengandalkan penyembunyian menu untuk permission.

5. **Kontrak permission/error**

   - Standarkan error `401`, `403`, `404`, `409`, dan validasi `422` untuk semua alur kolaborasi.
   - Tetapkan error code stabil agar FE tidak melakukan branching berdasarkan message bebas.

### P1 - Dibutuhkan untuk alur kolaborasi lengkap

1. **Akses Open/Private**

   - Field access type pada kolaborasi.
   - Open dapat menerima request bergabung.
   - Private tidak menampilkan CTA bergabung dan hanya dapat menerima undangan admin.
   - Request yang ditolak tidak dapat diajukan ulang sesuai aturan bisnis.

2. **Role anggota dan lifecycle anggota**

   - Daftar role standar, termasuk `Lainnya` dan keterangan tambahan.
   - Satu role utama dikunci setelah diterima.
   - Admin dapat mengeluarkan anggota.
   - Anggota dapat keluar dengan alasan.
   - Admin wajib menunjuk pengganti sebelum keluar.

3. **Workspace**

   - Resource workspace berdasarkan collaboration ID.
   - Chat dan file khusus workspace.
   - Permission hapus: pemilik konten atau admin.
   - Status `read_only` setelah kolaborasi selesai.
   - Endpoint upload, list, detail, dan delete file.

4. **Penyelesaian kolaborasi**

   - Endpoint complete collaboration.
   - Field hasil akhir, output file/link, dan visibility public/private.
   - Detail dan cover terkunci setelah completed.
   - Output publik dapat dibaca semua user; output privat hanya anggota.

5. **Feed kolaborasi**

   - Section `Kolaborasi Anda`, `Kolaborasi Selesai`, dan `Kolaborasi Berjalan`.
   - Filter topik global.
   - Sorting/ranking ditetapkan BE.
   - Pagination/infinite loading untuk carousel.
   - Empty state tetap dapat dibedakan FE berdasarkan section.

### P2 - Penyempurnaan produk

1. Event notifikasi untuk follow, request, approval/rejection, member baru, perubahan admin, upvote, reply, mention, dan perubahan penting.
2. Mention dan pengiriman notifikasi real-time atau polling yang disepakati.
3. Pagination dan sorting final pada diskusi, users, notifications, dan kolaborasi.
4. Audit log untuk perubahan permission, admin, status, dan output.
5. Upload media dengan batas ukuran, MIME type, storage, dan lifecycle yang terdokumentasi.

## Hal yang perlu dikonfirmasi FE sebelum kontrak dikunci

1. Apakah guest benar-benar tidak boleh membuka profile publik dan katalog kolaborasi, atau hanya tidak boleh memakai fitur aksi?
2. Apakah onboarding draft hanya disimpan di browser seperti Business Summary, atau perlu endpoint draft BE?
3. Nilai enum final untuk status kolaborasi: misalnya `ACTIVE`, `COMPLETED`, `READ_ONLY`.
4. Nilai enum final untuk access type: `OPEN` dan `PRIVATE`.
5. Format card kolaborasi yang diperlukan: nama field cover, description, topic, status, member count, dan CTA.
6. Apakah notifikasi harus real-time, atau polling sudah cukup untuk versi pertama?
7. Apakah workspace file menggunakan storage yang sudah ada atau perlu provider/storage baru?

## Catatan integrasi untuk FE saat ini

- Semua endpoint protected membutuhkan Bearer token.
- Jangan mengandalkan urutan field JSON; gunakan nama field dari OpenAPI.
- Upload asset umum menggunakan `multipart/form-data` ke `POST /api/v1/assets` dengan field `file` dan opsional `category`.
- Asset umum disimpan di Supabase pada pola `users/{userID}/{category}/{uuid}.{ext}`. FE menggunakan `url` yang dikembalikan BE; jangan membentuk URL storage sendiri.
- Attachment chat tetap memakai bucket `chat-attachments`; pesan chat direct/group tersimpan di Firebase pada pola `chats/direct/{conversationId}/messages/{messageId}` atau `chats/groups/{groupChatId}/messages/{messageId}`.
- PostgreSQL tetap menyimpan metadata chat untuk pagination, read/unread, permission, dan audit; Firebase digunakan untuk sinkronisasi realtime dan record chat deterministik.
- Gunakan `onboarding_status` dan `onboarding_completed` hanya untuk status dasar sampai modul Topics selesai.
- Pada diskusi, tampilkan field background, objective, expected_output, dan media_url bila tersedia.
- Untuk upvote, FE dapat memakai `/api/v1/discussions/:id/upvote` dan `DELETE` pada path yang sama; alias `/votes` tetap tersedia.
- Untuk request kolaborasi dari detail diskusi, gunakan `POST /api/v1/discussions/:discussionId/collaboration-requests` dengan body `proposedContribution`.
- Untuk kirim ulang verifikasi, FE dapat memakai `/api/v1/auth/resend-verification`; alias lama `/resend-verification-email` tetap tersedia.
- Daftar topik resmi dapat diambil dari `/api/v1/topics`; topik minat user dikelola melalui `/api/v1/profile/interests`.
- Jangan menampilkan tombol edit diskusi jika response menunjukkan kolaborasi aktif; enforcement final di BE masih perlu diselesaikan pada P0.
- Untuk notifikasi, endpoint yang tersedia adalah:
  - `GET /api/v1/notifications`
  - `GET /api/v1/notifications/unread-count`
  - `PATCH /api/v1/notifications/:id/read`
  - `POST /api/v1/notifications/:id/read`
  - `PATCH /api/v1/notifications/read-all`
  - `DELETE /api/v1/notifications/:id`
  - `DELETE /api/v1/notifications`
- Detail endpoint dan schema tetap merujuk ke `docs/openapi.yaml`.

## Kesimpulan

BE sudah cukup untuk integrasi awal auth, profile dasar, topics/interests dasar, feed diskusi, replies/vote, request kolaborasi dasar, asset upload, chat realtime, notifications dasar, dan publication feed. Implementasi berikutnya yang masih dibutuhkan adalah validasi onboarding final, approval admin topic, relasi diskusi-kolaborasi yang mengembalikan `collaboration_id`, permission guest, workspace file/read-only, lifecycle anggota, dan penyelesaian kolaborasi.

