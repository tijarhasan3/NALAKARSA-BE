# Organization API - FE Handoff

Dokumen ini menjelaskan penggunaan data instansi dan perusahaan dalam negeri.

## Konsep

Data organisasi disimpan pada tabel `organizations` dan terpisah dari tabel universitas `institutions`.

Contoh organisasi:

- Instansi pemerintah
- Perusahaan dalam negeri
- Lembaga riset
- BUMN

Field `type` tidak digunakan. FE hanya mengirim nama organisasi.

## Base URL

```text
{{url}}/api/v1
```

## 1. Mengambil daftar organisasi

Digunakan saat dropdown organisasi pertama kali dibuka.

```http
GET {{url}}/api/v1/organizations/search?page=1&limit=10
```

Header:

```http
Accept: application/json
```

Parameter query:

| Parameter | Tipe | Wajib | Default | Keterangan |
|---|---|---:|---:|---|
| `page` | integer | Tidak | `1` | Nomor halaman |
| `limit` | integer | Tidak | `10` | Jumlah data per halaman, maksimal `20` |
| `q` | string | Tidak | kosong | Kata pencarian nama organisasi |

Jika `q` tidak dikirim, backend mengembalikan organisasi aktif berdasarkan urutan nama.

## 2. Mencari organisasi

Saat user mengetik pada input organisasi, kirim query `q`.

```http
GET {{url}}/api/v1/organizations/search?q=Telkom&page=1&limit=10
```

Pencarian tidak membedakan huruf besar dan kecil.

Contoh:

```text
q=bank
```

Backend akan mengembalikan organisasi yang namanya mengandung `bank`.

## 3. Response pencarian

```json
{
  "data": [
    {
      "id": "uuid",
      "name": "PT Telkom Indonesia (Persero) Tbk",
      "countryCode": "ID",
      "country": "Indonesia"
    }
  ],
  "meta": {
    "current_page": 1,
    "total_pages": 1,
    "total_items": 1,
    "limit": 10
  },
  "message": "Organizations retrieved successfully"
}
```

## 4. Menambahkan organisasi baru

Gunakan endpoint ini jika hasil pencarian tidak menemukan organisasi yang dicari.

```http
POST {{url}}/api/v1/organizations
Content-Type: application/json
```

Request body:

```json
{
  "name": "PT Contoh Indonesia"
}
```

Field request:

| Field | Tipe | Wajib | Keterangan |
|---|---|---:|---|
| `name` | string | Ya | Nama organisasi, minimal 2 dan maksimal 255 karakter |

Tidak perlu mengirim field berikut:

```json
{
  "type": "company"
}
```

Field `type` sudah tidak digunakan.

## 5. Response tambah organisasi

Status berhasil: `201 Created`

```json
{
  "data": {
    "id": "uuid",
    "name": "PT Contoh Indonesia",
    "countryCode": "ID",
    "country": "Indonesia"
  },
  "message": "Organization created successfully"
}
```

## 6. Validasi nama

Backend otomatis:

- Menghapus spasi berlebih.
- Menghapus nomor penomoran di awal atau akhir.
- Menolak nama di bawah 2 karakter.
- Menolak nama di atas 255 karakter.
- Menolak karakter yang tidak didukung.
- Menolak kata yang masuk blacklist.
- Mencegah duplikasi nama melalui unique index database.

Contoh normalisasi:

```text
"  PT   Contoh Indonesia  "
```

Menjadi:

```text
"PT Contoh Indonesia"
```

## 7. Alur FE saat input organisasi

1. FE membuka dropdown organisasi.
2. FE memanggil `GET /organizations/search?page=1&limit=10` tanpa `q`.
3. FE menampilkan data pada dropdown.
4. User mengetik nama organisasi.
5. FE memanggil endpoint dengan parameter `q`.
6. Jika data ditemukan, user memilih salah satu organisasi.
7. Jika data tidak ditemukan, FE menampilkan pilihan `Tambah organisasi baru`.
8. FE memanggil `POST /organizations` dengan body hanya `name`.
9. FE menggunakan data response untuk mengisi nilai organisasi pada form.

## 8. Penggunaan pada profile atau registrasi

Nilai yang dikirim ke profile atau registrasi adalah nama organisasi:

```json
{
  "affiliation": "PT Contoh Indonesia"
}
```

Untuk update profile:

```http
PATCH {{url}}/api/v1/users/me
Authorization: Bearer {{accessToken}}
Content-Type: application/json
```

```json
{
  "affiliation": "PT Contoh Indonesia"
}
```

Backend menyimpan nilai tersebut pada kolom `users.affiliation`.

## 9. Error yang perlu ditangani FE

### Nama tidak valid

Status: `400 Bad Request`

```json
{
  "error": {
    "code": "BAD_REQUEST",
    "message": "name contains invalid characters or length"
  }
}
```

### Organisasi mengandung kata blacklist

Status: `400 Bad Request`

FE dapat menampilkan pesan bahwa nama organisasi tidak dapat digunakan.

## 10. Catatan implementasi

- Jangan mengirim `type`.
- Gunakan debounce saat memanggil pencarian, sekitar 300-500 ms.
- Jangan memanggil pencarian untuk setiap perubahan karakter tanpa debounce.
- Gunakan `meta.total_pages` untuk tombol halaman berikutnya.
- Saat response `data` kosong, tampilkan tombol tambah organisasi baru.
- Simpan nama organisasi, bukan hanya label tampilan sementara.
