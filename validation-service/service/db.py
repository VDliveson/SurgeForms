import os
from pymongo import MongoClient

# Retrieve the MongoDB connection string from the environment variable
MONGODB_URI = os.getenv("MONGODB_URI")


def connect_to_mongodb():
    try:
        client = MongoClient(MONGODB_URI)
        # db = client.your_database
        # collection = db.your_collection

        # # Perform MongoDB operations as needed
        # # For example, you can insert a document
        # document = {"key": "value"}
        # collection.insert_one(document)

        print("Connected to MongoDB successfully!")

        # Return the client and other objects if needed
        return client

    except Exception as e:
        print(f"Error connecting to MongoDB: {str(e)}")


if __name__ == "__main__":
    # Call the connect_to_mongodb function
    mongo_client, mongo_db, mongo_collection = connect_to_mongodb()

    # Example: Query all documents from the collection
    documents = mongo_collection.find()
    for document in documents:
        print(document)

    # Example: Close the MongoDB connection when done
    mongo_client.close()
