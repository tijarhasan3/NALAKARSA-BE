# FE: Vote dan Unvote Diskusi

## Ringkasan

User yang sudah login dapat memberikan vote pada diskusi dan membatalkan vote sebelumnya.

Gunakan access token pada header:

```http
Authorization: Bearer ACCESS_TOKEN
```

Jangan menggunakan `refresh_token` sebagai Bearer token.

## Vote Diskusi

```http
POST /api/v1/discussions/{discussionId}/votes
Authorization: Bearer ACCESS_TOKEN
```

Tidak membutuhkan request body.

Contoh:

```http
POST /api/v1/discussions/ea18ed37-4c08-4f22-ab2a-03900659bc1b/votes
```

Response sukses: HTTP `201 Created`.

```json
{
  "message": "Upvoted successfully",
  "data": {
    "id": "ea18ed37-4c08-4f22-ab2a-03900659bc1b",
    "upvote_count": 1,
    "has_upvoted": true
  }
}
```

Response mengembalikan detail diskusi. FE dapat mengambil nilai terbaru dari:

- `data.upvote_count`: jumlah vote
- `data.has_upvoted`: status vote user saat ini

## Unvote Diskusi

```http
DELETE /api/v1/discussions/{discussionId}/votes
Authorization: Bearer ACCESS_TOKEN
```

Tidak membutuhkan request body.

Response sukses: HTTP `200 OK`.

```json
{
  "message": "Upvote removed successfully",
  "data": {
    "id": "ea18ed37-4c08-4f22-ab2a-03900659bc1b",
    "upvote_count": 0,
    "has_upvoted": false
  }
}
```

## Implementasi FE

```javascript
async function toggleVote(discussionId, hasUpvoted) {
  const method = hasUpvoted ? "DELETE" : "POST";

  const response = await fetch(
    `/api/v1/discussions/${discussionId}/votes`,
    {
      method,
      headers: {
        Authorization: `Bearer ${accessToken}`
      }
    }
  );

  if (!response.ok) {
    throw new Error("Gagal memperbarui vote diskusi");
  }

  return response.json();
}
```

Setelah response diterima, update state dari response terbaru. Jangan hanya menambah atau mengurangi angka secara manual karena status vote dapat berubah akibat request lain.

## Error yang Perlu Ditangani

| HTTP | Kondisi | Pesan FE |
|---|---|---|
| `401` | Token tidak ada atau tidak valid | Silakan login kembali |
| `404` | Diskusi tidak ditemukan | Diskusi tidak tersedia |
| `409` | Vote sudah pernah diberikan | Diskusi sudah Anda vote |
| `409` | Unvote tanpa vote aktif | Vote tidak ditemukan |

## Catatan

- Tombol vote hanya dapat digunakan user yang sudah login.
- Gunakan `has_upvoted` untuk menentukan icon/status tombol.
- `POST` digunakan untuk vote.
- `DELETE` digunakan untuk unvote.
- `discussionId` harus berupa UUID diskusi yang valid.
