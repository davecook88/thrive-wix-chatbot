# Thrive Wix Chatbot

**A sophisticated, AI-powered chatbot for Wix websites, designed to provide a seamless and intelligent user experience.**

This project showcases a custom-built chatbot integrated into a Wix website. The chatbot leverages a powerful Go backend and a Web Component-based frontend to interact with users, retrieve data from Wix APIs, and personalize recommendations. The AI agent can suggest classes, manage student information in the CRM, and much more.

## ✨ Features

*   **🤖 AI-Powered Agent:** A smart assistant that understands user queries and provides intelligent responses.
*   **🌐 Wix Integration:** Seamlessly connects with Wix APIs to fetch real-time data on classes, courses, and appointments.
*   ** Personalized Recommendations:** The AI agent analyzes user preferences to suggest the most suitable classes and services.
*   **📝 CRM Integration:** Automatically saves and updates student information in the Wix CRM.
*   **⚡️ High-Performance Backend:** Built with Go, ensuring a fast, scalable, and reliable server.
*   **🎨 Custom Frontend:** A sleek and modern chat interface built with Web Components for maximum compatibility.
*   **🔒 Secure Authentication:** Protects user data and ensures secure communication between the frontend and backend.

## 🚀 Technologies Used

*   **Backend:** Go (Golang)
*   **Frontend:** Web Components (JavaScript, HTML, CSS)
*   **API:** Wix API
*   **Database:** Firebase
*   **Frameworks & Libraries:**
    *   Gin (Go web framework)
    *   Firebase Admin SDK for Go

## 🏗️ Architecture

The project follows a modern, decoupled architecture:

*   **Frontend:** A custom chat element built with Web Components is embedded into the Wix website. This component communicates with the backend via a secure WebSocket connection.
*   **Backend:** A Go server handles all the business logic, including:
    *   Managing the WebSocket connection with the frontend.
    *   Authenticating and authorizing users.
    *   Interacting with the Wix API to fetch and update data.
    *   Processing user messages with the AI agent.
    *   Storing and retrieving data from Firebase.

## 🏁 Getting Started

To get a local copy up and running, follow these simple steps.

### Prerequisites

*   Go 1.22.4 or higher
*   Node.js and npm
*   A Wix developer account
*   A Firebase project

### Installation

1.  **Clone the repo:**
    ```sh
    git clone https://github.com/your_username/thrive-wix-chatbot.git
    ```
2.  **Install backend dependencies:**
    ```sh
    cd thrive-go-server
    go mod download
    ```
3.  **Set up environment variables:**
    Create a `.env` file in the `thrive-go-server` directory and add the necessary environment variables for Wix and Firebase.
4.  **Run the backend server:**
    ```sh
    go run main.go
    ```

## 📄 License

Distributed under the MIT License. See `LICENSE` for more information.