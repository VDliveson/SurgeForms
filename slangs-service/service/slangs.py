import json
import pika
import base64
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
# Your Gmail user
GMAIL_USER = "vanshajduggal1234@gmail.com"


def load_credentials():
    creds = None
    creds = Credentials(
        token_uri="https://accounts.google.com/o/oauth2/token",
        client_id=os.getenv('CLIENT_ID'),
        client_secret=os.getenv('CLIENT_SECRET'),
        refresh_token=os.getenv('REFRESH_TOKEN'),
        token=os.getenv('ACCESS_TOKEN'),
    ) 
    if not creds or not creds.valid:
        if creds and creds.expired and creds.refresh_token:
            creds.refresh(Request())

    return creds


def get_slangs(message):
    city = None
    word = None

    for answer in message['message']['createdAnswers']:
        qs_id = answer['question']['_id']
        if(qs_id == "657eea8d9aaccdc0baf9cdfb"):
            city = answer['text']
        elif(qs_id == "657eea8d9aaccdc0baf9cdfc"):
            word = answer['text']
            
    # Fetch from database or API call
    answer = 'आनंद'
    return answer,word
  

def create_message(to, sender, subject, body):
    msg = EmailMessage()
    msg = MIMEText(body,'html')
    msg['To'] = to
    msg['From'] = sender
    msg['Subject'] = subject

    # msg.add_header('Content-Type','text/html')
    # msg.set_payload(body)
    encodedMsg = base64.urlsafe_b64encode(msg.as_bytes()).decode()
    return { 'raw': encodedMsg }
    # message = f"From: {sender}\nTo: {to}\nSubject: {subject}\n\n{body}"
    # return {'raw': base64.urlsafe_b64encode(message.encode()).decode()}

def send_message(service, sender, message):
    try:
        sent_message = service.users().messages().send(userId=sender, body=message).execute()
        return sent_message
    except Exception as error:
        LOGGER.error(f"An error occurred while sending the message: {error}")
    

def generate_response(data):
    try:
        answer,word = get_slangs(data)
        # answer = 2
        # word = 'sd'
        email_to = data['message']['metadata']['mail']
        user_id = data['message']['createdResponse']['user']
        response_id = data['message']['createdResponse']['_id']
        subject = "Your searched slangs".format(user_id)
        message = """
            <html>
            <body>
                <p>Hi user,</p>
                <p>Your query for slangs for word <b>{}</b> returned <b>{}</b></p>
            </body>
            </html>
            """.format(word,answer)
        

        credentials = load_credentials()

        service = build('gmail', 'v1', credentials=credentials)

        email_message = create_message(email_to, GMAIL_USER, subject, message)
        send_message(service, GMAIL_USER, email_message)
        LOGGER.info(f"Email sent successfully to {email_to} with subject '{subject}'")
    except Exception as e:
        LOGGER.error(f"Error sending email: {e}")

def main():
    pass
    # send_email(data)

if __name__ == '__main__':
    main()
