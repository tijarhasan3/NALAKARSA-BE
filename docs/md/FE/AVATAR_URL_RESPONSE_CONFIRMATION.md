# Konfirmasi Update FE–BE: Standarisasi `avatarUrl`

Tanggal: 19 November 2026  
Status: BE sudah disesuaikan

## Perubahan utama

BE sekarang menggunakan satu nama field avatar untuk response:

```json
"avatarUrl": "https://..."
```

Field lama berikut tidak lagi dikirim pada response:

```json
"avatar_url": "https://..."
```

## Endpoint yang terpengaruh

- `GET /api/v1/users/{id}`
- `GET /api/v1/users/{id}/followers`
- `GET /api/v1/users/{id}/following`
- `GET /api/v1/users`
- `POST /api/v1/users/me/avatar`
- `GET /api/v1/discussions`
- `GET /api/v1/discussions/{id}`

## Contoh response user

```json
{
  "id": "user-uuid",
  "fullName": "Tyezar Hasan Khusaeni",
  "avatarUrl": "https://teiqzolsyzqbbyxsnkmf.supabase.co/storage/v1/object/public/avatars/user-uuid.png"
}
```

## Contoh response followers/following

```json
{
  "id": "user-uuid",
  "fullName": "Nama User",
  "role": "profesional",
  "avatarUrl": "https://..."
}
```

## Penggunaan di FE

Gunakan langsung:

```js
user.avatarUrl
```

Untuk diskusi:

```js
discussion.avatarUrl
discussion.creator?.avatarUrl
```

Jangan lagi menggunakan:

```js
user.avatar_url
```

## Alur penyimpanan asset

Foto profil pada alur upload asset disimpan di bucket `assets` dengan struktur:

```text
assets/users/{user_id}/avatars/{file}
```

Foto diskusi menggunakan struktur:

```text
assets/users/{user_id}/discussion/{file}
```

Jika FE meng-upload foto profil melalui:

```http
POST /api/v1/assets
```

dengan:

```text
category=avatars
```

BE otomatis menyimpan URL public asset tersebut ke profil user. FE tidak perlu melakukan update avatar kedua kali.

## Catatan URL

BE menormalisasi path avatar lama menjadi URL Supabase lengkap. Contoh:

```text
/avatars/user.png
```

menjadi URL public dari bucket yang sesuai. Untuk foto profil:

```text
https://teiqzolsyzqbbyxsnkmf.supabase.co/storage/v1/object/public/assets/users/{user_id}/avatars/user.png
```

Untuk file lama yang sudah terlanjur berada di `assets/users/{user_id}/avatars` tetapi belum tercatat di database, BE akan mencari file terbaru saat profil pertama kali dibaca, menyimpan URL-nya ke `users.avatar_url`, lalu mengembalikan `avatarUrl`. Jika tetap kosong, berarti file avatar memang belum ditemukan di folder tersebut.

