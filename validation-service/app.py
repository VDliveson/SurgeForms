import amqpstorm
import json
import time
import os
from dotenv import load_dotenv

load_dotenv()
from service.mail import send_email
from service.logger import LOGGER

RABBITMQ = os.getenv("RABBITMQ") or "amqp://localhost:5672"

service_id = "validation"
exchange = "daisy1"
queue = "mail_queue"


def connect_queue():
    connected = False
    while not connected:
        try:
            connection = amqpstorm.UriConnection(RABBITMQ)
            channel = connection.channel()

            # Declare exchange and queue
            channel.exchange.declare(
                exchange=exchange, exchange_type="direct", durable=True
            )
            LOGGER.info(f"Connected to RabbitMQ exchange {exchange}")

            channel.queue.declare(queue=queue, durable=True)
            channel.queue.bind(exchange=exchange, queue=queue, routing_key=service_id)

            def callback(message):
                content = message.body
                LOGGER.info(f" [x] Received message: {content}\n")
                message = json.loads(content)
                try:
                    send_email(message)
                except Exception as e:
                    LOGGER.error(f"Failed to send email: {e}")
                    return  # Exit the callback function
                LOGGER.info("Work completed\n")

            channel.basic.qos(prefetch_count=2)
            channel.basic.consume(queue=queue, callback=callback, no_ack=True)
            LOGGER.info(" [*] Waiting for messages. To exit, press Ctrl+C")

            connected = True
            try:
                channel.start_consuming()
            except KeyboardInterrupt:
                LOGGER.info("Interrupted, closing connection.")
                channel.stop_consuming()
                connection.close()
                LOGGER.info("Connection closed.")
                return

        except Exception as e:
            LOGGER.error(e)
            LOGGER.error("Failed to connect to RabbitMQ. Retrying in 5 seconds.")
            time.sleep(5)


if __name__ == "__main__":
    connect_queue()
