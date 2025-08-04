# up-cli

📦 `up-cli` is a cross-platform command-line tool to upload media files to cloud storage providers such as **Supabase**, **Cloudinary**, and **Backblaze**. It also stores metadata in a **Neon PostgreSQL** database.

## 🚀 Features

- Upload files to:
  - 📤 [Supabase Storage](https://supabase.com/docs/guides/storage)
  - ☁️ [Cloudinary](https://cloudinary.com/)
  - 💾 [Backblaze B2](https://www.backblaze.com/b2/cloud-storage.html)
- Metadata persistence in [Neon PostgreSQL](https://neon.tech/)
- Interactive provider selection
- Built-in version command (`up-cli --version`)
- Cross-platform binary builds and installation scripts

## 🛠️ Stack

- **Golang** for CLI
- **Supabase**, **Cloudinary**, **Backblaze B2** for storage
- **NeonDB** for database
- **Cobra** for CLI command management

## 📦 Installation

You can use the [installation script](./install.sh) to install the binary:

```bash
curl -sSf https://raw.githubusercontent.com/smokeyshawn18/up-cli/main/install.sh | bash


🔧 Usage

up-cli upload ./your-file.png
```
