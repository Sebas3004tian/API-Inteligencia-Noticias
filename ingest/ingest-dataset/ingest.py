from datasets import load_dataset
import requests
import uuid
import time

API_ENDPOINT = "http://localhost:8080/index"
BATCH_SIZE = 1000

def main():
    print("Descargando dataset...")

    ds = load_dataset(
        "csv",
        data_files="https://huggingface.co/datasets/MarcOrfilaCarreras/spanish-news/resolve/main/data.csv",
    )["train"]

    print("Columnas:", ds.column_names)
    print("Filas:", len(ds))

    batch = []

    for i, row in enumerate(ds):
        content = row.get("text", "") if isinstance(row, dict) else row

        item = {
            "id": str(uuid.uuid4()),
            "title": "",
            "description": "",
            "content": content,
        }

        batch.append(item)

        if len(batch) >= BATCH_SIZE:
            send_batch(batch)
            batch = []

    if batch:
        send_batch(batch)

    print("Proceso completado.")

def send_batch(batch):
    print(f"Enviando batch con {len(batch)} artículos...")
    try:
        resp = requests.post(API_ENDPOINT, json=batch)
        if resp.status_code not in (200, 201):
            print(f" Error batch: {resp.status_code} -> {resp.text}")
        else:
            print("✓ Batch enviado")
    except Exception as e:
        print(f" Error enviando batch: {e}")

    time.sleep(30)

if __name__ == "__main__":
    main()
