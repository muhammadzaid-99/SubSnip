# vidsub-ocr

A distributed OCR-based subtitle extraction system built with Go, RabbitMQ, Docker, and Python. Easily scale workers and extract subtitles from videos using videocr (PaddleOCR).

## ⚙️ Architecture

* **Producer**: A Go server that accepts subtitle extraction tasks.
* **Workers**: Python-based consumers that process tasks from RabbitMQ.
* **Scalable**: Set `replicas` in `docker-compose.yml` to deploy multiple workers concurrently.
* **OCR Engine**: [videocr](https://github.com/devmaxxing/videocr-PaddleOCR) for accurate subtitle detection.

## 🚀 Features

* 🔁 Distributed task queue via RabbitMQ
* 📥 Video download + (commented) subtitle extraction logic
* 🧠 Configurable parameters: language, frame skip, time range, bounding box
* 📊 Track task status: `queued`, `processing`, `done`
* 🐳 Dockerized and horizontally scalable

## 🛠️ Requirements

* Docker + Docker Compose
* Go 1.21+
* Python 3.10+
* `requirements.txt` includes both GPU and CPU versions of PaddleOCR
  → **Uncomment only one** based on your system.

## 🔧 Getting Started

```bash
docker-compose up --build
```

Submit tasks to the producer server. Workers will pick them up automatically.

## 📝 Notes

* Subtitle extraction logic in Python is currently commented for stability. 
* You can scale worker replicas using:

```yaml
services:
  worker:
    deploy:
      replicas: 3
```

## TODO

* Re-enable and test subtitle extraction logic
* Add frontend for task tracking
* Integrate results export

