#!/usr/bin/env python3
import os
import sys
import socket

def main():
    target_ip = os.environ.get("RSCTF_TARGET_IP")
    target_port = int(os.environ.get("RSCTF_TARGET_PORT", "9000"))
    round_num = os.environ.get("RSCTF_ROUND")
    flag = os.environ.get("RSCTF_FLAG")

    if not target_ip:
        print("Missing target IP", file=sys.stderr)
        sys.exit(3)  # InternalError

    # 1. Cek Konektivitas (Offline check)
    s = socket.socket(socket.AF_INET, socket.SOCK_STREAM)
    s.settimeout(5.0)
    try:
        s.connect((target_ip, target_port))
    except (socket.timeout, ConnectionRefusedError, OSError) as e:
        print(f"Connection failed: {e}", file=sys.stderr)
        sys.exit(2)  # Offline

    # 2. Cek Fungsi Service / Verifikasi Flag (Mumble / OK check)
    try:
        banner = s.recv(1024).decode(errors="ignore")
        if "Welcome to Secure Vault" not in banner:
            print("Invalid banner received", file=sys.stderr)
            sys.exit(1)  # Mumble

        # Kirim request cek fungsi
        s.sendall(b"PING\n")
        res = s.recv(1024).decode(errors="ignore")
        if "PONG" not in res:
            print("Service logic ping failed", file=sys.stderr)
            sys.exit(1)  # Mumble

        s.close()
        # Jika semua flow lolos:
        sys.exit(0)  # OK
    except socket.timeout:
        sys.exit(2)  # Offline
    except Exception as e:
        print(f"Error during check: {e}", file=sys.stderr)
        sys.exit(1)  # Mumble

if __name__ == "__main__":
    main()
