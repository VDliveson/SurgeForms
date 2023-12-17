import amqpstorm
import json
import time
import os
from dotenv import load_dotenv
from models.db import connect_to_mongodb 
import threading
from service.sheets import process_data

load_dotenv()

service_id = "sheets"
exchange = "daisy1"
messages = []
mongo_client = None
RABBITMQ = os.getenv("RABBITMQ") or "amqp://localhost:5672"

def connect_queue():
    connected = False
    while not connected:
        try:
            connection = amqpstorm.UriConnection(
                RABBITMQ
            )
            channel = connection.channel()

            channel.exchange.declare(
                exchange=exchange, exchange_type="direct", durable=True
            )
            print(f"Connected to RabbitMQ exchange {exchange}")

            queue = "sheets_queue"
            channel.queue.declare(queue=queue, durable=True)
            channel.queue.bind(exchange=exchange, queue=queue, routing_key=service_id)

            def callback(message):
                content = message.body
                print(f" [x] Received message: {content}")
                message = json.loads(content)
                # messages.append(message)
                
                process_data(message, mongo_client)

                print("Work completed")

            channel.basic.qos(prefetch_count=2)
            channel.basic.consume(queue=queue, callback=callback, no_ack=True)
            print(" [*] Waiting for messages. To exit, press Ctrl+C")
            connected = True
            channel.start_consuming()

        except Exception as e:
            print(e)
            print(f"Failed to connect to RabbitMQ. Retrying in 5 seconds.")
            time.sleep(5)

if __name__ == "__main__":
    mongo_client = connect_to_mongodb()
    connect_queue()
