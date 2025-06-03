# Nerdify Audiobook

A web audiobook application with Go backend and React frontend that provides free audiobook streaming from LibriVox.

## 🚀 **Setup and Installation**

### **1. Clone Repository**
```bash
git clone 
cd nerdify-audiobook/SOA
```

### **2. Setup PostgreSQL Database**

#### **Create Database:**
CREATE DATABASE audiobook;

### **3. Setup Environment Variables**

Create your .env file, adjust it according to your PostgreSQL configuration:

```env
DB_HOST=localhost
DB_PORT=5432
DB_NAME=audiobook
DB_USER=postgres
DB_PASSWORD=postgre_password_here
GOOGLE_CLIENT_ID=your_google_client_id (kosongin aja)
GOOGLE_CLIENT_SECRET=your_google_client_secret (kosongin aja)
SERVER_HOST=http://localhost:8000
```

### **4. Install Backend Dependencies**

```bash
# Install Go dependencies
go mod download
go mod tidy
```


### **5. Seed Database**

Run the seeder to create tables and initial data:

```bash
go run seed.go
```

### **7. Setup Frontend**

```bash
cd frontend

# Install dependencies
yarn install
```

### **8. Start Backend Server**

```bash
# Run backend server
go run main.go book.go
```

### **9. Start Frontend Development Server**

Open a new terminal:

```bash
cd frontend
yarn start
```

