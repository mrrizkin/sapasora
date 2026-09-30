# GoFiber + Inertia.js + Vite

This is a simple starter template for a GoFiber + Inertia.js + Vite application.

## Features

- GoFiber
- Inertia.js
- Vite
- Tailwind CSS
- TypeScript
- ESLint
- Prettier

## Getting Started

1. Clone the repository
2. Run `pnpm install` to install dependencies
3. Run `pnpm dev` to start the development server
4. Open http://localhost:3000 to view the application

## Docker (development / build from source)

Docker akan membangun TDLib native, frontend, dan aplikasi Go secara otomatis.

```bash
cp .env.example .env
# isi TELEGRAM_API_ID dan TELEGRAM_API_HASH di .env

docker compose up --build
```

Aplikasi tersedia di http://localhost:3000. Data PostgreSQL, sesi TDLib, dan storage disimpan dalam Docker volumes.

Untuk menghentikan container:

```bash
docker compose down
```

## Production deploy (pakai image yang sudah di-build CI)

Setiap push ke `main` otomatis di-build dan di-publish oleh GitHub Actions ke
GitHub Container Registry sebagai `ghcr.io/mrrizkin/sapasora`, dengan tag versi
bergaya CalVer `vYYYY.WW.PATCH` (mis. `v2026.39.0`) sekaligus tag `latest`.

Di server, cukup pull image tersebut tanpa perlu build TDLib/Go/Vite sama sekali:

```bash
cp .env.example .env
# isi .env sesuai kebutuhan production (lihat docs/PRODUCTION_READINESS.md)

# opsional: pin ke versi tertentu, default-nya "latest"
echo "IMAGE_TAG=v2026.39.0" >> .env

docker compose -f docker-compose.prod.yml pull
docker compose -f docker-compose.prod.yml up -d
```

Untuk update ke versi baru, ganti `IMAGE_TAG` di `.env` lalu ulangi `pull` + `up -d`.
Lihat tag yang tersedia di halaman Packages repo ini di GitHub.

## API documentation (OpenAPI)

Spec OpenAPI 3.0 (JSON + YAML) di-generate otomatis saat build image (via
`toolbox swagger`), tapi endpoint `/docs/*` sengaja diblokir saat
`APP_ENV=production` supaya peta API tidak terbuka ke publik. Untuk mengambil
spec-nya (misal untuk dikirim ke client yang mau implementasi API), extract
langsung dari image yang sudah di-publish tanpa perlu build atau menjalankan
app:

```bash
./scripts/extract-openapi.sh              # image :latest -> ./openapi-export
./scripts/extract-openapi.sh v2026.40.0    # pin ke versi tertentu
```

Hasilnya `openapi-export/openapi.json` dan `openapi.yaml`, siap di-import ke
Postman/Insomnia atau dirender pakai Swagger UI/Redoc.

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.
