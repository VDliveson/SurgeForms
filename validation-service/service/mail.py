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

# Your Gmail API credentials (downloaded from the Google Cloud Console)
credentials_path = r"D:\Atlan project\Atlan-assignment\validation-service\service\oauth-token.json"

# Google API Scopes required for sending emails
SCOPES = ["https://www.googleapis.com/auth/gmail.send"]

# Your Gmail user
GMAIL_USER = "vanshajduggal1234@gmail.com"

# token_path = r'D:\Atlan project\Atlan-assignment\validation-service\service\oauth-token.json'
# f = open(token_path)
# data = json.load(f)


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
        # else:
        #     flow = InstalledAppFlow.from_client_secrets_file(credentials_path, SCOPES)
        #     creds = flow.run_local_server(port=0)

        # with open(token_path, 'w') as token:
        #     token.write(creds.to_json())

    return creds


def validate_message(message):
    qs_ans_dict = {}
    income = 0
    savings = 0
    for answer in message['message']['createdAnswers']:
        qs_id = answer['question']['_id']
        if(qs_id == "6577ee2d868435b880cf243b"):
            savings = int(answer['text'])
        if(qs_id == "6577ee2d868435b880cf243c"):
            income = int(answer['text'])
        
    
    if(savings > income):
        return True
    else:
        return False
  

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
        print(f"An error occurred while sending the message: {error}")
    

def send_email(data):
    try:
        if(validate_message(data)):            
            email_to = data['message']['metadata']['mail']
            user_id = data['message']['createdResponse']['user']
            response_id = data['message']['createdResponse']['_id']
            subject = "Flagging user {} for invalid response".format(user_id)
            message = """
                <html>
                <body>
                    <p>Hi collector,</p>
                    <p>Invalid response id: <b>{}</b> from user id: <b>{}</b> for salary and income.</p>
                </body>
                </html>
                """.format(response_id,user_id)
            

            credentials = load_credentials()

            service = build('gmail', 'v1', credentials=credentials)

            email_message = create_message(email_to, GMAIL_USER, subject, message)
            send_message(service, GMAIL_USER, email_message)
            print(f"Email sent successfully to {email_to} with subject '{subject}'")
    except Exception as e:
        print(f"Error sending email: {e}")

def main():
    pass
    # send_email(data)

if __name__ == '__main__':
    main()
