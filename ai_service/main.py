from fastapi import FastAPI, UploadFile, File, HTTPException
from pydantic import BaseModel
from PIL import Image
import pytesseract
import io
from ultralytics import YOLO
from collections import Counter

app = FastAPI()

# Load Model (will download if not present, but Dockerfile caches it)
try:
    model = YOLO("yolov8n.pt") 
except Exception as e:
    print(f"Error loading YOLO model: {e}")
    model = None

class AnalysisResult(BaseModel):
    items: list
    raw_text: str

import pika
import json
import threading
import time
import os

def process_job(ch, method, properties, body):
    print(f" [x] Received job: {body}")
    # In real app: Download image, process, then hit a webhook or update DB
    # For now we just log it
    
    try:
        data = json.loads(body)
        print(f"Processing image for user: {data.get('user_id')}")
        # Logic from analyze_image would go here
    except Exception as e:
        print(f"Error processing job: {e}")

    ch.basic_ack(delivery_tag=method.delivery_tag)

def start_consumer():
    rabbitmq_host = os.getenv("RABBITMQ_HOST", "rabbitmq")
    while True:
        try:
            print("Connecting to RabbitMQ...")
            connection = pika.BlockingConnection(
                pika.ConnectionParameters(host=rabbitmq_host))
            channel = connection.channel()

            channel.queue_declare(queue='image_processing_queue', durable=True)

            channel.basic_qos(prefetch_count=1)
            channel.basic_consume(queue='image_processing_queue', on_message_callback=process_job)

            print(' [*] Waiting for messages. To exit press CTRL+C')
            channel.start_consuming()
        except Exception as e:
            print(f"Connection failed: {e}. Retrying in 5s...")
            time.sleep(5)

# Start consumer in background thread
threading.Thread(target=start_consumer, daemon=True).start()

@app.get("/health")
def health_check():
    return {"status": "ok", "service": "ai-engine-real"}

@app.post("/analyze", response_model=AnalysisResult)
async def analyze_image(file: UploadFile = File(...)):
    if not model:
        raise HTTPException(status_code=500, detail="AI Model not loaded")

    # 1. Read Image
    try:
        image_data = await file.read()
        image = Image.open(io.BytesIO(image_data))
    except Exception:
        raise HTTPException(status_code=400, detail="Invalid image file")

    # 2. Object Detection (YOLO)
    # Run inference
    results = model(image)
    
    detected_items = []
    # Count occurrences
    labels = []
    for result in results:
        for box in result.boxes:
            class_id = int(box.cls[0])
            label = model.names[class_id]
            labels.append(label)

    counts = Counter(labels)
    for label, count in counts.items():
        detected_items.append({
            "name": label,
            "qty": count,
            "confidence": 1.0 # Simplified for aggregation
        })

    # 3. OCR (Text Extraction)
    try:
        raw_text = pytesseract.image_to_string(image)
    except Exception as e:
        raw_text = f"Error performing OCR: {e}"

    return {
        "items": detected_items,
        "raw_text": raw_text.strip()
    }

if __name__ == "__main__":
    import uvicorn
    uvicorn.run(app, host="0.0.0.0", port=5000)
