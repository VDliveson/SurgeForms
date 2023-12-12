import pika
import json
import time
import os
import asyncio
from dotenv import load_dotenv
load_dotenv()
from service.mail import send_email 
RABBITMQ_PORT = os.getenv("RABBITMQ_PORT")

service_id = "validation"
exchange = "daisy1"
messages = []


def connect_queue():
    try:
        # connection = pika.BlockingConnection(pika.ConnectionParameters('amqp://localhost:5672'))  # Update with your RabbitMQ server details
        connection = pika.BlockingConnection(
            pika.ConnectionParameters(host="localhost", port=RABBITMQ_PORT)
        )  #
        channel = connection.channel()

        channel.exchange_declare(
            exchange=exchange, exchange_type="direct", durable=True
        )
        print(f"Connected to RabbitMQ exchange {exchange}")

        queue = "verification_queue"
        channel.queue_declare(queue=queue, durable=True)
        channel.queue_bind(exchange=exchange, queue=queue, routing_key=service_id)

        def callback(ch, method, properties, body):
            content = body.decode("utf-8")
            print(f" [x] Received message: {content} from ID '{id}'")
            message = json.loads(content)
            messages.append(message)
            
            # print(message['message']['metadata'])
            send_email(message)
            # sheets_create_and_add(message['message'])

            # Uncomment the following block if you want to send acknowledgment after processing
            # acknowledgment = {"message": "Microservice1 acknowledged message"}
            # ch.basic_publish(
            #     exchange='',
            #     routing_key=properties.reply_to,
            #     properties=pika.BasicProperties(
            #         correlation_id=properties.correlation_id
            #     ),
            #     body=json.dumps(acknowledgment)
            # )

            print("Work completed")
            # ch.basic_ack(delivery_tag=method.delivery_tag)
        channel.basic_qos(prefetch_count=2)
        channel.basic_consume(queue=queue, on_message_callback=callback, auto_ack=True)
        print(" [*] Waiting for messages. To exit, press Ctrl+C")
        channel.start_consuming()

    except Exception as e:
        print(str(e))


if __name__ == "__main__":
    connect_queue()
