import os
from dotenv import load_dotenv
load_dotenv()

db = None

from datetime import datetime
from bson import ObjectId
from google.oauth2.credentials import Credentials
from google_auth_oauthlib.flow import InstalledAppFlow
from google.auth.transport.requests import Request
from googleapiclient.discovery import build


# Google API Scopes required for sending emails
SCOPES = ['https://www.googleapis.com/auth/spreadsheets']


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


def process_data(data,mongo_client):
    try:
        global db 
        db = mongo_client['test']
        response = data["message"]["createdResponse"]["_id"]
        form = data["message"]["createdResponse"]["form"]["_id"]
        title = data["message"]["createdResponse"]["form"]["title"]
        user =  data["message"]["createdResponse"]["user"]
        createdAnswer = data["message"]["createdAnswers"]
        answers_array = []
        for answer in createdAnswer:
            temp = {}
            temp["question_id"] = answer["question"]["_id"]
            temp["question_text"] = answer["question"]["text"]
            temp["answer_text"] = answer["text"];  
            answers_array.append(temp)          
        
        insert_response(form, response, answers_array)
        add_data_to_sheet(form,answers_array)

    except Exception as e:
        print('Error adding data to sheet:', e)

def insert_response(form, response, answers_array):
    try:        
        answers = []
        for answer in answers_array:
            res = {
                "response": ObjectId(response),
                "form": ObjectId(form),
                "question": ObjectId(answer["question_id"]), 
                "text": answer["answer_text"],
                "createdAt": datetime.today().replace(microsecond=0)
            }
            answers.append(res)
        db.responses.insert_many(answers)
        print("Answers inserted successfully into the database")

    except Exception as e:
        print(f"Error inserting responses into the database: {e}")
    


def add_data_to_sheet(form, answers):
    try:
        values = []
        for value in answers:
            values.append(value['answer_text'])
        rows = values

        resource = {
            'values': [rows]
        }
        
        sheet_id = get_form_sheet(form)

        spreadsheet_id = sheet_id
        range_ = 'Sheet1'
        sheets = build('sheets', 'v4', credentials=load_credentials())
        response = sheets.spreadsheets().values().append(
            spreadsheetId=spreadsheet_id,
            range=range_,
            valueInputOption='RAW',
            body=resource
        ).execute()

        print('Data inserted successfully into spreadsheet')
    except Exception as e:
        print('Unable to insert data into the database',e)
        raise Exception(e)
        



def get_form_sheet(form_id):
    try:
        existing_entry = db.sheets.find_one({'form': form_id})
        if existing_entry:
            print("Sheet entry already exists in MongoDB.")
            return existing_entry['sheet']  
              
        spreadsheet_title = 'Form' + "_" + form_id[-5:]
        sheets = build('sheets', 'v4', credentials=load_credentials())
        spreadsheet = sheets.spreadsheets().create(
            body={'properties': {'title': spreadsheet_title}}
        ).execute()

        spreadsheet_id = spreadsheet['spreadsheetId']
        print(f"Spreadsheet ID: {spreadsheet_id}")

        # Insert document into MongoDB
        sheet = {
            '_id': ObjectId(),
            'form': form_id,
            'sheet': spreadsheet_id,
            'createdAt': datetime.today().replace(microsecond=0)
        }
        db.sheets.insert_one(sheet)

        print("Sheet information inserted into MongoDB successfully.")

        return spreadsheet_id

    except Exception as e:        
        print(f"Error creating Google Sheet: {e}")
        raise Exception(e)
