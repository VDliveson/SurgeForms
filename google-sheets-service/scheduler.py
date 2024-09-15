from models.db import connect_to_mongodb
import os


mongo_client = connect_to_mongodb()
from dotenv import load_dotenv

load_dotenv()

from datetime import datetime
from bson import ObjectId
from google.oauth2.credentials import Credentials
from google_auth_oauthlib.flow import InstalledAppFlow
from google.auth.transport.requests import Request
from googleapiclient.discovery import build


form = "657716650200915e6885fe76"


def load_credentials():
    creds = None
    creds = Credentials(
        token_uri="https://accounts.google.com/o/oauth2/token",
        client_id=os.getenv("CLIENT_ID"),
        client_secret=os.getenv("CLIENT_SECRET"),
        refresh_token=os.getenv("REFRESH_TOKEN"),
        token=os.getenv("ACCESS_TOKEN"),
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


def insert_data_to_sheets(form):
    try:
        db = mongo_client["test"]
        responses = db.responses.find({"form": ObjectId(form)})
        all_responses = []
        sheet_id = db.sheets.find_one({"form": ObjectId(form)})["sheet"]
        print(sheet_id)

        print(f"Responses for Form ID: {form}")
        for response in responses:
            curr = []
            for answer in response["answers"]:
                text = answer["answer_text"]
                curr.append(text)
            all_responses.append(curr)
        if len(all_responses):
            add_data_to_sheet(form, all_responses, sheet_id)
            db.responses.delete_many({"form": ObjectId(form)})
            print(f"Responses saved to sheets and deleted succesfully : {form}")
        print("Process completed successfully")

    except Exception as e:
        print(f"Error fetching responses: {e}")


def add_data_to_sheet(form, answers, sheet_id):
    try:
        resource = {"values": answers}

        spreadsheet_id = sheet_id
        range_ = "Sheet1"
        sheets = build("sheets", "v4", credentials=load_credentials())
        response = (
            sheets.spreadsheets()
            .values()
            .append(
                spreadsheetId=spreadsheet_id,
                range=range_,
                valueInputOption="RAW",
                body=resource,
            )
            .execute()
        )

        print("Data inserted successfully into spreadsheet")
    except Exception as e:
        print("Unable to insert data into the database", e)
        raise Exception(e)


if __name__ == "__main__":
    insert_data_to_sheets(form)
