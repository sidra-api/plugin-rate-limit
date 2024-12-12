# Plugin Rate Limit

## Deskripsi
Plugin Rate Limit digunakan untuk membatasi jumlah permintaan yang dapat dilakukan oleh sebuah IP dalam satu menit pada Sidra Api. Plugin ini membantu melindungi layanan backend dari serangan DoS (Denial of Service) atau penyalahgunaan oleh klien.

---

## Cara Kerja
1. **Pemeriksaan Header**
   - Plugin memeriksa header `X-Real-Ip` untuk mendapatkan alamat IP klien.
   - Jika header ini tidak ditemukan, permintaan akan ditolak dengan status `400 Bad Request`.

2. **Pembatasan Permintaan**
   - Plugin mencatat jumlah permintaan yang dilakukan oleh setiap IP dalam satu menit menggunakan peta (map).
   - Batas default adalah **5 permintaan per menit**.

3. **Respon**
   - Jika jumlah permintaan masih di bawah batas:
     - Status: `200 OK`
     - Body: "Request allowed"
   - Jika batas permintaan terlampaui:
     - Status: `429 Too Many Requests`
     - Body: "Rate limit exceeded"

---

## Konfigurasi
- **Batas Permintaan**
  - Dapat dikonfigurasi langsung pada file `main.go`:
    ```go
    const rateLimitPerMinute = 5
    ```

---

## Cara Menjalankan
1. Pastikan Anda sudah menginstal **Sidra Api**.
2. Tambahkan plugin ini ke direktori `plugins/rate-limit/main.go` pada Sidra Api.
3. Kompilasi dan jalankan Sidra Api.
4. Plugin akan otomatis terhubung melalui UNIX socket pada path `/tmp/ratelimit.sock`.

---

## Pengujian

### Endpoint
- **URL**: Endpoint mana saja yang dikonfigurasi untuk melewati plugin Rate Limit.

### Langkah Pengujian
1. Kirim beberapa permintaan dari IP yang sama menggunakan alat pengujian seperti Postman atau curl.
2. Respons yang diharapkan:
   - Jika permintaan masih dalam batas:
     - Status: `200 OK`
     - Body: "Request allowed"
   - Jika batas terlampaui:
     - Status: `429 Too Many Requests`
     - Body: "Rate limit exceeded"

---

## Catatan Penting
- **Header IP**: Pastikan header `X-Real-Ip` dikirim oleh proxy atau load balancer sebelum mencapai Sidra Gateway.
- **Pengaturan Ulang**: Plugin akan mengatur ulang peta permintaan setiap satu menit menggunakan goroutine.

---

## Lisensi
MIT License
