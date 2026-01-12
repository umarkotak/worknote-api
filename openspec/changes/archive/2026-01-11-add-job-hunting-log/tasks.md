## 1. Database Layer

- [x] 1.1 Create migration `002_create_job_applications.sql` with job_applications table
- [x] 1.2 Create migration `003_create_job_application_logs.sql` with job_application_logs table

## 2. Model Layer

- [x] 2.1 Add `JobApplication` struct to `model/model.go`
- [x] 2.2 Add `JobApplicationLog` struct to `model/model.go`

## 3. Contract Layer

- [x] 3.1 Add job application request/response types to `contract/contract.go`
- [x] 3.2 Add job application log request/response types to `contract/contract.go`

## 4. Repository Layer

- [x] 4.1 Create `repos/job_application_repo/` with Initialize, prepared statements
- [x] 4.2 Implement Create, GetByID, GetByUserID (with search/filter), Update, Delete
- [x] 4.3 Create `repos/job_application_log_repo/` with Initialize, prepared statements
- [x] 4.4 Implement Create, GetByID, GetByJobApplicationID, Update, Delete

## 5. Service Layer

- [x] 5.1 Create `services/job_application_service/` with business logic
- [x] 5.2 Implement CRUD operations for job applications
- [x] 5.3 Implement CRUD operations for job application logs
- [x] 5.4 Add authorization checks (user can only access own applications)

## 6. Handler Layer

- [x] 6.1 Create `handlers/job_application_handler/` with HTTP handlers
- [x] 6.2 Implement POST/GET/PUT/DELETE for job applications
- [x] 6.3 Implement GET with search/filter query params
- [x] 6.4 Implement POST/GET/PUT/DELETE for job application logs (nested)

## 7. Routing

- [x] 7.1 Register job application routes in `main.go` (protected with auth middleware)
- [x] 7.2 Initialize new repositories in `main.go`

## 8. Validation

- [x] 8.1 Build compiles successfully (`go build ./...`)
- [ ] 8.2 Manually test all endpoints via curl/Postman
- [ ] 8.3 Verify authorization (users can only see their own data)
