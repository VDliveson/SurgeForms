from flask import Flask, jsonify
from threading import Thread
from app import connect_queue
import os
from dotenv import load_dotenv

load_dotenv()

app = Flask(__name__)


def start_queue_consumer():
    try:
        connect_queue()
    except KeyboardInterrupt:
        # Graceful shutdown on Ctrl+C
        print("Received KeyboardInterrupt, stopping consumer...")


@app.route("/")
def index():
    return jsonify({"status": "Flask app is running!"})


if __name__ == "__main__":
    # Start the RabbitMQ consumer in a separate thread
    consumer_thread = Thread(target=start_queue_consumer)
    consumer_thread.start()

    # Run the Flask app with debug mode
    app.run(debug=True, port=os.getenv("PORT"), use_reloader=True)

    # Wait for the consumer thread to finish
    consumer_thread.join()