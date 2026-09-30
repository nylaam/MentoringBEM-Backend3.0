# RSCTF Attack & Defense (A&D) Problemset Template

Repository ini adalah template dan panduan resmi untuk struktur challenge berformat **Attack & Defense (A&D)** di platform RSCTF.

---

## 📁 1. Struktur Folder Repository Challenge

```text
my-ctf-challenges/
├── .gzevent                                 # (Wajib di root repo) Manifest Event & Konfigurasi A&D
└── challenges/
    └── AD/
        ├── Pwn/
        │   └── secure-vault/
        │       ├── challenge.yml            # Manifest spesifik challenge
        │       ├── dist/                    # File yang dibagikan ke peserta (attachment/handout)
        │       │   └── vault.zip
        │       ├── src/                     # Source code vulnbox / service container
        │       │   ├── Dockerfile           # Dockerfile service (jika di-build oleh RSCTF)
        │       │   ├── Makefile
        │       │   └── server.c
        │       └── checker/                 # SLA Checker (Python sandbox)
        │           ├── run.py               # Entry point checker script
        │           └── requirements.txt     # Dependency Python (wajib exact pin name==version)
        └── Web/
            └── banking-api/
                ├── challenge.yml
                ├── dist/
                │   └── src.zip
                ├── src/
                │   ├── Dockerfile
                │   ├── app.py
                │   └── requirements.txt
                └── checker/
                    ├── run.py
                    └── requirements.txt
```

---

## ⚙️ 2. Penjelasan Komponen & Format File

### A. Event Manifest: `.gzevent` (Root Repo)
File ini diletakkan di root repository untuk mendefinisikan event/game baru dan konfigurasi round A&D.

```yaml
title: "RSCTF 2026 Attack & Defense"
summary: "Annual National CTF Competition"
start: "2026-10-01T08:00:00Z"
end: "2026-10-01T16:00:00Z"
hidden: true
acceptWithoutReview: true
teamMemberCountLimit: 4
ad:
  tickSeconds: 120              # Durasi 1 round / tick dalam detik (misal 2 menit)
  flagLifetimeTicks: 3          # Masa berlaku flag (misal 3 round)
  warmupSeconds: 300            # Masa pemanasan sebelum round 1 (misal 5 menit)
  resetCooldownMinutes: 10      # Cooldown jika tim mereset vulnbox mereka
  allowSnapshotDownload: false  # Izinkan tim mendownload snapshot vulnbox
```

---

### B. Challenge Manifest: `challenge.yml`
Diletakkan di dalam folder masing-masing challenge (misal: `challenges/AD/Pwn/secure-vault/challenge.yml`).

```yaml
name: "Secure Vault"
author: "dimas"
type: "AttackDefense"           # Wajib "AttackDefense"
category: "Pwn"
description: |
  Layani request vault dan amankan data nasabah.
  Endpoint service: `nc <team-ip> 9000`

# Konfigurasi container service vulnbox peserta
container:
  exposePort: 9000              # Port service yang dibuka ke jaringan A&D
  memoryLimit: 512              # Batas RAM dalam MB
  cpuCount: 1                   # Batas CPU core
  storageLimit: 1024            # Batas Storage dalam MB
  # Jika containerImage dikosongkan, RSCTF akan otomatis mem-build Dockerfile dari folder src/
  # containerImage: "ghcr.io/myorg/secure-vault:latest"

# Konfigurasi spesifik Attack & Defense
ad:
  allowEgress: false            # false = blokir koneksi keluar dari vulnbox (rekomendasi keamanan)
  allowSelfReset: true          # Izinkan peserta mereset instance vulnbox mereka via dashboard
  sshRequiresFlag: false        # Akses SSH tidak memerlukan flag
  selfHosted: false             # false = di-host di infrastruktur platform (bukan BYOC)

# File yang disediakan untuk di-download peserta
provide: "dist/vault.zip"
```

---

### C. Service Source: `src/Dockerfile`
Tempat source code vulnbox yang akan di-deploy untuk setiap tim.

```dockerfile
FROM ubuntu:24.04

RUN apt-get update && apt-get install -y --no-install-recommends \
    socat \
    gcc \
    libc6-dev \
    && rm -rf /var/lib/apt/lists/*

WORKDIR /app
COPY src/ /app/
RUN gcc -O2 -fstack-protector server.c -o /app/server

EXPOSE 9000

# Service dijalankan via daemon/socat
CMD ["socat", "TCP-LISTEN:9000,reuseaddr,fork", "EXEC:/app/server"]
```

---

### D. SLA Checker: `checker/run.py` & `requirements.txt`
RSCTF menjalankan SLA checker menggunakan isolated Python venv sandbox (bukan Docker container terpisah) untuk efisiensi dan keamanan.

#### 1. Kontrak Environment Variables (`RSCTF_*`)
Saat RSCTF mengecek service suatu tim di setiap round, RSCTF mengeksekusi `run.py` dengan environment variables berikut:

- `RSCTF_ACTION`: Aksi checker (default `"check"`).
- `RSCTF_TARGET_IP`: Alamat IP vulnbox tim target.
- `RSCTF_TARGET_PORT`: Port service target (misal `9000` atau `5000`).
- `RSCTF_ROUND`: Nomor round saat ini (misal `5`).
- `RSCTF_TEAM_ID`: ID tim target.
- `RSCTF_CHALLENGE_ID`: ID challenge.
- `RSCTF_FLAG`: Flag unik yang harus ada/ditanam pada round tersebut.

#### 2. Kontrak Exit Code

| Exit Code | Status di RSCTF | Keterangan |
| :--- | :--- | :--- |
| **0** | **OK** | Service berjalan normal, SLA lolos, flag berhasil diverifikasi/ditemukan. |
| **1** | **Mumble** | Service merespons tetapi data korup, fitur rusak, atau flag tidak sesuai. |
| **2** | **Offline** | Port tertutup, connection refused, atau timeout (down). |
| **3** | **InternalError** | Script checker error / crash logic / missing IP. |

#### 3. Aturan `checker/requirements.txt`
RSCTF mewajibkan format `package==x.y.z`. Tidak diperbolehkan menggunakan format range (`>=`, `~=`) atau local build hooks:
```text
pwntools==4.14.0
requests==2.32.3
```
