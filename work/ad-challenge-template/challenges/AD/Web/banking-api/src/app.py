import os
from pathlib import Path
from flask import Flask, jsonify, request

app = Flask(__name__)

FLAG_FILE = Path(os.environ.get("RSCTF_FLAG_FILE", "/flag"))

ACCOUNTS = {
    "acc_1001": {"name": "Alice", "balance": 1500, "history": ["Initial deposit: $1500"]},
    "acc_1002": {"name": "Bob", "balance": 2400, "history": ["Initial deposit: $2400"]},
}

def get_current_flag() -> str:
    """Membaca flag ronde berjalan secara real-time dari filesystem."""
    try:
        return FLAG_FILE.read_text(encoding="utf-8").strip()
    except Exception:
        return "flag{not-delivered-yet}"

@app.route("/health", methods=["GET"])
def health():
    """Endpoint untuk verifikasi kesehatan service oleh SLA Checker."""
    return "ok", 200

@app.route("/api/balance", methods=["GET"])
def get_balance():
    account_id = request.args.get("account_id", "acc_1001")
    if account_id not in ACCOUNTS:
        return jsonify({"status": "error", "message": "Account not found"}), 404
    return jsonify({
        "status": "success",
        "account_id": account_id,
        "name": ACCOUNTS[account_id]["name"],
        "balance": ACCOUNTS[account_id]["balance"],
    }), 200

@app.route("/api/transfer", methods=["POST"])
def transfer():
    data = request.get_json(silent=True) or {}
    from_acc = data.get("from_acc")
    to_acc = data.get("to_acc")
    amount = data.get("amount")

    if not from_acc or not to_acc or amount is None:
        return jsonify({"status": "error", "message": "Missing transfer parameters"}), 400

    try:
        amount = int(amount)
        if amount <= 0:
            return jsonify({"status": "error", "message": "Amount must be positive"}), 400
    except ValueError:
        return jsonify({"status": "error", "message": "Invalid amount"}), 400

    if from_acc not in ACCOUNTS or to_acc not in ACCOUNTS:
        return jsonify({"status": "error", "message": "Account does not exist"}), 404

    if ACCOUNTS[from_acc]["balance"] < amount:
        return jsonify({"status": "error", "message": "Insufficient balance"}), 400

    ACCOUNTS[from_acc]["balance"] -= amount
    ACCOUNTS[to_acc]["balance"] += amount
    ACCOUNTS[from_acc]["history"].append(f"Sent ${amount} to {to_acc}")
    ACCOUNTS[to_acc]["history"].append(f"Received ${amount} from {from_acc}")

    return jsonify({
        "status": "success",
        "message": "Transfer successful",
        "from_balance": ACCOUNTS[from_acc]["balance"],
    }), 200

@app.route("/api/statement", methods=["GET"])
def get_statement():
    """
    Endpoint vuln: Membocorkan flag ronde jika parameter secret=1
    atau jika account_id bernilai 'admin' / 'vault'.
    """
    account_id = request.args.get("account_id", "acc_1001")
    current_flag = get_current_flag()

    if account_id in ("admin", "vault") or request.args.get("secret") == "1":
        return jsonify({
            "status": "success",
            "account_id": account_id,
            "statement": f"Confidential transaction data: {current_flag}",
            "flag": current_flag,
        }), 200

    if account_id not in ACCOUNTS:
        return jsonify({"status": "error", "message": "Account not found"}), 404

    return jsonify({
        "status": "success",
        "account_id": account_id,
        "history": ACCOUNTS[account_id]["history"],
    }), 200

if __name__ == "__main__":
    app.run(host="0.0.0.0", port=5000)
