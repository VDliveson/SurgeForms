# SurgeForms

### Overview
This is a project aimed at creating microservices for form collection and performing asynchronous operations on the form data that scales well across a large number of users.

### Tools and technologies used
- Docker
- Python
- RabbitMQ
- NodeJS
- External APIs

### Architecture Diagram
![image](https://github.com/user-attachments/assets/3ced4c7a-56ed-4e70-9ec4-947cafcb3a37)

<br><br>
> [!Important]  
> Some features are not implemented yet as shown in the diagram

### Running the application

### Running the Application

1. **Install Docker** :
   Ensure Docker is installed on your system. You can download it from [Docker's official website](https://www.docker.com/get-started).

2. **Clone the Repository** :
   Clone the repository using either of the following methods -
   - **HTTP**: `git clone https://github.com/VDliveson/SurgeForms`
   - **SSH**: `git clone git@github.com:VDliveson/SurgeForms.git`

3. **Add Environment Files** :
   Add a `.env` file in each microservice folder within the repository. Ensure all required environment variables are specified in these files as given in `.example.env` files

4. **Start the Application** :
   In the root directory of the repository, run the following command to start the microservices:
   ```bash
   docker compose up
   ```
> [!NOTE]  
> This may take a couple of minutes to initialize

5. **Stop the Application** :
   To terminate the application, run the following command -
   ```bash
   docker compose down
   ```
6. **Clean Up Docker Images ( Optional )** :
   To clean up unused Docker images and free up space, run the following command :
   ```bash
   docker image prune -a
   ```
   
