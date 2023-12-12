import os
from dotenv import load_dotenv
load_dotenv()
from twilio.rest import Client
import logging


account_sid = os.environ['TWILIO_ACCOUNT_SID']
auth_token = os.environ['TWILIO_AUTH_TOKEN']
client = Client(account_sid, auth_token)
from_number = os.environ['TWILIO_FROM_NUMBER']

def remove(string):
    return "".join(string.split())

def send_message(user_id,response_id,title,to_number):
    try:
        body = "Hi ! Your account with user_id {} has received form response id: {} for form titled : {}".format(user_id, response_id,title)
        message = client.messages \
                    .create(
                        body=body,
                        from_= from_number,
                        to=remove(to_number)
                    )

        print("Sent message with sid :",message.sid)
    except Exception as e:
        raise Exception("Failed to send message")

def process_message(data):
    try:
        response_id = data['message']['createdResponse']['_id']
        user_id = data['message']['createdResponse']['user']
        phone = data['message']['metadata']['phone number']
        form_title = data['message']['createdResponse']['form']['title']
        send_message(user_id,response_id,form_title,phone)
    except Exception as e:
        logging.error(e)
