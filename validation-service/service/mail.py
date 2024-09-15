import json
import base64
import smtplib
import os
from dotenv import load_dotenv

load_dotenv()
from email.mime.multipart import MIMEMultipart
from email.mime.text import MIMEText
from email.mime.base import MIMEBase

from google.oauth2.credentials import Credentials
from google_auth_oauthlib.flow import InstalledAppFlow
from email.message import EmailMessage
from google.auth.transport.requests import Request

from googleapiclient.discovery import build


from .logger import LOGGER

SMTP_SERVER = os.getenv("SMTP_SERVER", "smtp.gmail.com")
SMTP_PORT = int(os.getenv("SMTP_PORT", 587))
SMTP_USER = os.getenv("SMTP_USER")
SMTP_PASSWORD = os.getenv("SMTP_PASSWORD")


def validate_message(message):
    income = 0
    savings = 0
    try:
        for answer in message["message"]["createdAnswers"]:
            qs_id = answer["question"]["_id"]
            if qs_id == "66bf6c93286c3acef094728a":
                savings = int(answer["text"])
            if qs_id == "66bf6c93286c3acef094728b":
                income = int(answer["text"])
    except Exception as e:
        LOGGER.error(f"Error parsing message: {e}")
        return False

    if savings > income:
        return True
    else:
        return False


def create_message(to, sender, subject, body):
    msg = MIMEMultipart()
    msg["To"] = to
    msg["From"] = sender
    msg["Subject"] = subject
    msg.attach(MIMEText(body, "html"))
    return msg


def send_message(smtp_server, smtp_port, smtp_user, smtp_password, message):
    try:
        with smtplib.SMTP(smtp_server, smtp_port) as server:
            server.starttls()  # Upgrade the connection to a secure encrypted SSL/TLS connection
            server.login(smtp_user, smtp_password)
            server.send_message(message)
        LOGGER.info(
            f"Email sent successfully to {message['To']} with subject '{message['Subject']}'"
        )
    except Exception as e:
        raise e


def send_email(data):
    try:
        if validate_message(data):
            email_to = data["message"]["metadata"]["mail"]
            user_id = data["message"]["createdResponse"]["user"]
            response_id = data["message"]["createdResponse"]["_id"]
            subject = f"Flagging user {user_id} for invalid response"
            message_body = f"""
                <html>
                <body>
                    <p>Hi collector,</p>
                    <p>Invalid response id: <b>{response_id}</b> from user id: <b>{user_id}</b> for salary and income.</p>
                </body>
                </html>
            """

            message = create_message(email_to, SMTP_USER, subject, message_body)
            LOGGER.info(f"Sending email to {email_to}")
            send_message(SMTP_SERVER, SMTP_PORT, SMTP_USER, SMTP_PASSWORD, message)
    except Exception as e:
        raise e
