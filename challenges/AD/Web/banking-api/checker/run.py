#!/usr/bin/env python3
import os
import sys
import requests

def main():
    target_ip = os.environ.get("RSCTF_TARGET_IP")
    target_port = int(os.environ.get("RSCTF_TARGET_PORT", "5000"))
    round_num = os.environ.get("RSCTF_ROUND")
    flag = os.environ.get("RSCTF_FLAG")

    if not target_ip:
        print("Missing target IP", file=sys.stderr)
        sys.exit(3)  # InternalError

    base_url = f"http://{target_ip}:{target_port}"

    # 1. Cek Konektivitas & Endpoint /health (Offline check)
    try:
        res = requests.get(f"{base_url}/health", timeout=5.0)
    except (requests.exceptions.ConnectTimeout, requests.exceptions.ConnectionError, requests.exceptions.Timeout) as e:
        print(f"Connection failed: {e}", file=sys.stderr)
        sys.exit(2)  # Offline
    except Exception as e:
        print(f"Unexpected error: {e}", file=sys.stderr)
        sys.exit(3)  # InternalError

    if res.status_code != 200 or res.text.strip() != "ok":
        print(f"Health check failed (status: {res.status_code}, text: {res.text})", file=sys.stderr)
        sys.exit(1)  # Mumble

    # 2. Cek Fungsi Service: GET /api/balance
    try:
        res = requests.get(f"{base_url}/api/balance?account_id=acc_1001", timeout=5.0)
        if res.status_code != 200:
            print(f"Balance check status invalid: {res.status_code}", file=sys.stderr)
            sys.exit(1)  # Mumble
        data = res.json()
        if data.get("status") != "success" or "balance" not in data:
            print("Balance check JSON payload invalid", file=sys.stderr)
            sys.exit(1)  # Mumble
    except (requests.exceptions.ConnectTimeout, requests.exceptions.ConnectionError, requests.exceptions.Timeout):
        sys.exit(2)  # Offline
    except Exception as e:
        print(f"Balance check error: {e}", file=sys.stderr)
        sys.exit(1)  # Mumble

    # 3. Cek Fungsi Service: POST /api/transfer
    try:
        res = requests.post(
            f"{base_url}/api/transfer",
            json={"from_acc": "acc_1001", "to_acc": "acc_1002", "amount": 10},
            timeout=5.0
        )
        if res.status_code != 200:
            print(f"Transfer check status invalid: {res.status_code}", file=sys.stderr)
            sys.exit(1)  # Mumble
        data = res.json()
        if data.get("status") != "success":
            print("Transfer failed according to API response", file=sys.stderr)
            sys.exit(1)  # Mumble
    except (requests.exceptions.ConnectTimeout, requests.exceptions.ConnectionError, requests.exceptions.Timeout):
        sys.exit(2)  # Offline
    except Exception as e:
        print(f"Transfer check error: {e}", file=sys.stderr)
        sys.exit(1)  # Mumble

    # 4. Cek Verifikasi Flag (jika RSCTF_FLAG disediakan)
    if flag:
        try:
            res = requests.get(f"{base_url}/api/statement?account_id=admin", timeout=5.0)
            if res.status_code != 200:
                print(f"Statement check status invalid: {res.status_code}", file=sys.stderr)
                sys.exit(1)  # Mumble
            data = res.json()
            if flag not in data.get("statement", "") and data.get("flag") != flag:
                print("Flag verification failed: flag not found in service response", file=sys.stderr)
                sys.exit(1)  # Mumble
        except (requests.exceptions.ConnectTimeout, requests.exceptions.ConnectionError, requests.exceptions.Timeout):
            sys.exit(2)  # Offline
        except Exception as e:
            print(f"Flag verification error: {e}", file=sys.stderr)
            sys.exit(1)  # Mumble

    # Semua check lolos
    sys.exit(0)

if __name__ == "__main__":
    main()
