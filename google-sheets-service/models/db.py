import os
from pymongo import MongoClient

# Retrieve the MongoDB connection string from the environment variable
MONGODB_URI = os.getenv("MONGODB_URI")


def connect_to_mongodb():
    try:
        client = MongoClient(MONGODB_URI)
        print("Connected to MongoDB successfully!")

        # Return the client and other objects if needed
        return client

    except Exception as e:
        print(f"Error connecting to MongoDB: {str(e)}")
        connect_to_mongodb()


if __name__ == "__main__":
    # Call the connect_to_mongodb function
    mongo_client= connect_to_mongodb()