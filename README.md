# Nihongowa

Nihongowa is an innovative chat application designed to simulate conversations in Japanese, akin to the functionality of WhatsApp. By integrating OpenAI's powerful language models, users can engage in dialogues in Japanese, receive translations, and thereby enhance their language learning experience. This application is perfect for language enthusiasts aiming to improve their Japanese conversational skills in a dynamic and interactive environment.

## Features

- Chat with OpenAI in Japanese.
- View translations of conversations to aid learning.
- Mimics a real chat application environment.
- Local and cloud deployment options.

## Getting Started

These instructions will get you a copy of the project up and running on your local machine for development and testing purposes. You can choose to run the application locally, through Docker, or deploy it to AWS using the AWS CDK.

### Prerequisites

Before you begin, ensure you have the following installed:

- Go (for running the server locally)
- Docker
- AWS CLI and AWS CDK (optional, for deployment to AWS)
- Android Studio (for generating the APK)

### Running Locally

1. **Clone the Repository**

```bash
git clone https://github.com/arymoraes/nihongowa.git
cd nihongowa
```

2. **Run the Server**

Add your OPENAI_API_KEY to the .env file in the api directory. Then run the server.

```bash
cd api/internal/server
go run server.go
```

3. **Initialize the Database**

Generate and run the Docker container for the database.

```bash
./startup.sh
```

4. **Run the Client**

Open the client directory in Android Studio and run the application on an emulator or physical device.
