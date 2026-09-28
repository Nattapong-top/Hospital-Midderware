# Hospital-Middleware

ระบบตัวกลางเชื่อมต่อข้อมูลโรงพยาบาล (Hospital Middleware System) พัฒนาด้วยภาษา Go โดยใช้สถาปัตยกรรม Clean Architecture / Domain-Driven Design (DDD) รองรับการจัดการข้อมูลพนักงาน, ระบบยืนยันตัวตนด้วย JWT, และการค้นหาข้อมูลผู้ป่วยข้ามโรงพยาบาล

---

## 1. ภาพรวมโปรเจกต์ (Overview)
- **Tech Stack:** Go (Golang), Gin Framework, PostgreSQL, Docker & Docker Compose, JWT Authentication, Bcrypt Password Hashing.
- **Architecture:** Clean Architecture (Domain, Application, Infrastructure, Delivery / HTTP Adapters).
---

## 2. วิธีติดตั้งและรันระบบ (Installation & Running)

### 2.1 รันผ่าน Docker Compose (แนะนำ)
โปรเจกต์นี้มี `docker-compose.yml` ที่จัดการทั้ง PostgreSQL, Go Application และ Nginx Reverse Proxy ครบถ้วนในตัว

1. โคลนโปรเจกต์และเข้าไปที่ไดเรกทอรีโปรเจกต์:
   ```bash
   git clone <repository-url>
   cd Hospital-Middleware
   ```
2. รันบริการทั้งหมดด้วย Docker Compose (พร้อม Build Image ใหม่):
   ```bash
   docker-compose up --build
   ```
   หรือรันแบบ Background (Detached mode):
   ```bash
   docker-compose up --build -d
   ```
3. คำสั่ง Docker Compose ที่มีประโยชน์:
   - ดู Logs: `docker-compose logs -f`
   - หยุดการทำงาน: `docker-compose down`
   - หยุดและลบข้อมูล Volume (รีเซ็ตฐานข้อมูล): `docker-compose down -v`

4. เข้าใช้งานผ่าน Nginx บนพอร์ต `80` (เช่น `http://localhost/staff/login`)

### 2.2 รันแบบ Local Development
1. ติดตั้ง Go (1.21+) และ PostgreSQL ในเครื่อง
2. สร้างฐานข้อมูล PostgreSQL และรัน Migration scripts ในโฟลเดอร์ `migrations/` ตามลำดับ:
   ```bash
   psql -h localhost -U postgres -d hospital_middleware -f migrations/000001_create_hospital_and_staff_tables.up.sql
   psql -h localhost -U postgres -d hospital_middleware -f migrations/000002_create_patients_table.up.sql
   psql -h localhost -U postgres -d hospital_middleware -f migrations/000003_add_version_to_staffs.up.sql
   ```
3. กำหนดตัวแปรสภาพแวดล้อม (Environment Variables) หรือไฟล์ `.env`
4. รันแอปพลิเคชัน:
   ```bash
   go run cmd/main.go
   ```

---

## 3. วิธีเรียกใช้งาน API (API Usage)

1. **เข้าสู่ระบบพนักงาน (Staff Login):**
   - **Method/Path:** `POST /staff/login`
   - **Body:** `{"username": "...", "password": "..."}`
   - **Response:** คืนค่า `token` (JWT) สำหรับนำไปใช้ยืนยันตัวตน

2. **ลงทะเบียนพนักงาน (Create Staff):**
   - **Method/Path:** `POST /staff/create`
   - **Body:** `{"username": "...", "password": "...", "hospital": "..."}`

3. **ค้นหาข้อมูลผู้ป่วย (Search Patient - Protected):**
   - **Method/Path:** `GET /patient/search?national_id=...`
   - **Headers:** `Authorization: Bearer <JWT_TOKEN>`

---

## 4. วิธีทดสอบระบบ (Testing)

รัน Unit Tests และ Integration Tests พร้อมตรวจสอบ Race Condition และ Code Coverage:
```bash
go test -v -cover ./...
```
