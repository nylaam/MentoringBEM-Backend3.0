FROM python:3.12-alpine

WORKDIR /app

COPY requirements.txt .
RUN pip install --no-cache-dir -r requirements.txt

COPY app.py .

# Path file flag default yang diperbarui otomatis oleh platform setiap round
ENV RSCTF_FLAG_FILE=/flag
RUN echo "flag{warmup-initial-flag-placeholder}" > /flag && chmod 0666 /flag

EXPOSE 5000

CMD ["gunicorn", "-w", "2", "-b", "0.0.0.0:5000", "app:app"]
