# **Product Management Backend**

## **Key Points**

### **Tech Stack:**
- **Go**: Gin framework for REST APIs.
- **PostgreSQL**: For database storage.
- **RabbitMQ**: For message queuing.
- **Redis**: For caching.
- **Cloudinary**: For image storage and compression.

---

### **Core Features:**

- **Product Management**:  
  CRUD operations for products, utilizing PostgreSQL for data storage.

- **Asynchronous Image Processing**:  
  - Images are sent to RabbitMQ upon product creation.  
  - A worker service processes images (compression via Cloudinary) and updates the database.

- **Caching**:  
  Redis is used to cache product data, reducing database load and improving performance.

---

### **Workflow:**

#### **Image Processing Workflow:**
1. **POST /products**:  
   - Saves product data in the database.  
   - Sends image URLs to a RabbitMQ queue.  
2. A **consumer service**:  
   - Retrieves messages from RabbitMQ.  
   - Compresses images using Cloudinary.  
   - Updates the `compressed_product_images` field in the database.

#### **GET with Caching:**
1. Product data is first fetched from Redis.  
2. If not cached, data is fetched from PostgreSQL and then cached for future requests.

---

### **Structured Logging:**
- Logs all API requests, errors, and processing events using **zap** for structured and efficient logging.

---

### **Testing:**
- Comprehensive unit and integration tests with **90%+ code coverage**.

---

### **Setup:**

1. **Environment Variables**:  
   Configure `.env` with:  
   - `DATABASE_URL`: PostgreSQL connection URL.  
   - `RABBITMQ_URL`: RabbitMQ connection URL.  
   - `CLOUDINARY_URL`: Cloudinary API URL.  
   - `REDIS_URL`: Redis connection URL.  

2. **Installation**:  
   ```bash
   git clone <repository-url>
   cd <repository-folder>
   go mod tidy
   go run main.go
   go run cmd/worker/main.go
# Output Screen Shot

### Post Request
![Screenshot 2024-12-10 162507](https://github.com/user-attachments/assets/b925d91d-a23b-4f39-b3f1-f0e797485534)
### Cloudinary Url for image
![Screenshot 2024-12-10 162524](https://github.com/user-attachments/assets/64f2ad75-79d9-463c-b9bb-59899626e07f)

### Before Redic Output
![Screenshot 2024-12-10 162549](https://github.com/user-attachments/assets/419b80b7-f2df-413b-8199-15f35ab9e240)
### Get Request
![Screenshot 2024-12-10 162624](https://github.com/user-attachments/assets/96558773-df26-4f64-b407-ff33ffe31528)

### After Get Request Redic Output
![Screenshot 2024-12-10 162648](https://github.com/user-attachments/assets/7f40378b-1f6e-4398-895c-0c9b513fc5ae)

### Database Product Table
![Screenshot 2024-12-10 164327](https://github.com/user-attachments/assets/f8c40e23-e463-4a1f-864e-271b88a0a241)
### Database User Table
![Screenshot 2024-12-10 164345](https://github.com/user-attachments/assets/76b9c416-c05a-43b8-a87a-879e7c292e9b)




