Recruitment Management System
This repository contains the backend server implementation for a Recruitment Management System. Below are the details and instructions for the project:

Task
Create a backend server for a Recruitment Management System with the following requirements:

Users can create their profile and upload Resumes (only PDF and DOCX formats allowed).

Admin users can create job openings.

Admin users can view all uploaded resumes and extracted data of applicants.

Applicants can view job openings and apply to them.

APIs
POST /signup: Create a user profile.
POST /login: Authenticate users and return a JWT token.
POST /uploadResume: Upload resume files (PDF or DOCX) for applicants.
POST /admin/job: Create job openings.
GET /admin/job/{job_id}: Fetch information about a job opening.
GET /admin/applicants: Fetch a list of all users in the system.
GET /admin/applicant/{applicant_id}: Fetch extracted data of an applicant.
GET /jobs: Fetch job openings.
GET /jobs/apply?job_id={job_id}: Apply to a particular job.
Models
User: Contains user details such as name, email, address, user type, password hash, and profile headline.
Profile: Contains applicant profile details such as resume file address, skills, education, experience, name, email, and phone.
Job: Contains job details such as title, description, posted date, total applications, company name, and posted by.
Third-party API
API Endpoint: https://api.apilayer.com/resume_parser/upload
Request Type: POST
Headers: "Content-Type": application/octet-stream, "apikey": this_is_the_api_key
