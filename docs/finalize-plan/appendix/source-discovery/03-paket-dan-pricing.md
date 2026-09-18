# 03 — Paket dan Pricing

> Harga dicatat dari toggle pricing di homepage Fonnte saat discovery. Periksa ulang sebelum pembelian.

## Text Only — bulanan (IDR)

| Paket | Kuota | Multics Agent | Harga |
|---|---:|---:|---:|
| Free | 1.000 pesan/bulan | 0 | Rp0 |
| Lite | 1.000 pesan/bulan | 0 | Rp25.000 |
| Regular | 10.000 pesan/bulan | 2 | Rp66.000 |
| Regular Pro | 25.000 pesan/bulan | 2 | Rp110.000 |
| Master | Unlimited | 4 | Rp175.000 |

## All Feature — bulanan (IDR)

| Paket | Kuota | Multics Agent | Harga |
|---|---:|---:|---:|
| Super | 10.000 pesan/bulan | 2 | Rp165.000 |
| Advanced | 25.000 pesan/bulan | 2 | Rp255.000 |
| Ultra | Unlimited | 4 | Rp355.000 |

## Text Only — tahunan (IDR)

| Paket | Kuota | Multics Agent | Harga |
|---|---:|---:|---:|
| Free | 1.000 pesan/bulan | 0 | Rp0 |
| Lite | 1.000 pesan/bulan | 0 | Rp250.000 |
| Regular | 10.000 pesan/bulan | 2 | Rp660.000 |
| Regular Pro | 25.000 pesan/bulan | 2 | Rp1.100.000 |
| Master | Unlimited | 4 | Rp1.750.000 |

## All Feature — tahunan (IDR)

| Paket | Kuota | Multics Agent | Harga |
|---|---:|---:|---:|
| Super | 10.000 pesan/bulan | 2 | Rp1.650.000 |
| Advanced | 25.000 pesan/bulan | 2 | Rp2.550.000 |
| Ultra | Unlimited | 4 | Rp3.550.000 |

## Harga efektif paket tahunan

Jika angka tahunan dibagi 12, perkiraan biaya per bulan adalah:

| Paket | Harga tahunan | Efektif/bulan | vs bulanan |
|---|---:|---:|---:|
| Lite | Rp250.000 | ~Rp20.833 | ~16,7% lebih murah |
| Regular | Rp660.000 | Rp55.000 | ~16,7% lebih murah |
| Regular Pro | Rp1.100.000 | ~Rp91.667 | ~16,7% lebih murah |
| Master | Rp1.750.000 | ~Rp145.833 | ~16,7% lebih murah |
| Super | Rp1.650.000 | Rp137.500 | ~16,7% lebih murah |
| Advanced | Rp2.550.000 | Rp212.500 | ~16,7% lebih murah |
| Ultra | Rp3.550.000 | ~Rp295.833 | ~16,7% lebih murah |

Perhitungan ini hanya membagi harga yang ditampilkan website, bukan konfirmasi mekanisme billing.

## Harga USD yang terlihat

Toggle USD menampilkan harga bulanan/tahunan berikut:

- Bulanan: Lite `$2.5`, Regular `$6`, Regular Pro `$10`, Master `$14`, Super `$13`, Advanced `$22`, Ultra `$29`.
- Tahunan: Lite `$25`, Regular `$60`, Regular Pro `$100`, Master `$140`, Super `$130`, Advanced `$220`, Ultra `$290`.
- Free tetap `$0`.

## Feature framing

Homepage membagi paket menjadi **Text Only** dan **All Feature**. Daftar label fitur yang ditampilkan meliputi:

- Personal/group send
- Text, schedule, recurring, template
- Button *(deprecated)*
- Attachment
- Autoreply/autoreply spreadsheet
- Webhook
- API
- Remove watermark
- Device notification
- Multics Agent

Perbedaan visual feature inclusion perlu dibaca bersama state ikon/check pada halaman pricing; snapshot teks tidak cukup untuk memetakan setiap entitlement secara pasti. Dokumentasi API mengonfirmasi bahwa attachment melalui `url`/`file` dibatasi ke Super/Advanced/Ultra.

## Aturan paket dari dokumentasi

- Setiap device memiliki paket sendiri; paket tidak dibagi antar-device.
- Paket tidak dapat upgrade/downgrade secara prorata; paket saat ini akan terminate saat upgrade/downgrade.
- Bulan dihitung 30 hari, tahun 365 hari.
- Subscribe ulang paket yang sama menambahkan quota dan masa berlaku.
- Upgrade/downgrade mengakhiri quota dan expiration paket sebelumnya.

## Analisis product/monetization

**Bagus:** free tier dengan kuota yang sama secara praktis untuk development menurunkan risiko trial; tangga kuota jelas; harga Indonesia sangat accessible.

**Potensi perbaikan:**

1. Jelaskan secara eksplisit perbedaan Text Only vs All Feature dengan comparison table, bukan hanya dua kelompok kartu.
2. Tampilkan harga per bulan efektif untuk paket tahunan dan besaran penghematan.
3. Jelaskan “message quota” ketika satu pesan dikirim ke banyak target.
4. Beri warning yang lebih menonjol tentang kebijakan terminasi quota saat upgrade/downgrade.
5. Tandai semua fitur deprecated/limited tepat di kartu paket.
6. Sinkronkan batas attachment antara dokumen (10 MB) dan contoh error API (di bawah 4 MB).

## Sumber

- <https://fonnte.com/>
- <https://docs.fonnte.com/about-package/>
- <https://docs.fonnte.com/api-send-message/>
