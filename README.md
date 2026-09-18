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

## Docker

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

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.
