# First Step

I’ve built an AI-based platform called First Step for generating security reports on any software.

I believe that the very first step before installing any software should be verifying its security and trustworthiness.

The platform consists of two parts: an API and a web client that interacts with it. The web interface provides a user-friendly way to generate security reports. Users simply provide an input: a software title, company name, website link, or even a brief description if the exact title is forgotten.

One AI agent identifies the software product the user wants to check and passes that information to the second AI agent — the Security Assessor.

The Security Assessor can search Google and visit various websites, gathering detailed information on vendor reputation, known vulnerabilities, and other real-world data relevant to software security.

Once the report is ready, the user sees a comprehensive overview of the software, including vendor details, key security issues, an overall security score, and a verdict on whether the software is safe to use — and under what conditions.

Reports are cached on the server and are reproducible. By design, all reports remain accessible via a unique identifier. Even when a report becomes outdated and a newer version is generated, the original report is always available.

## User perspective
The use case is simple: a user visits the website to check whether a given software is secure.
They enter the software information into an input field and click “Get report.”
It takes around 30 seconds to generate a report, which then becomes instantly accessible to anyone.
The user receives clear and meaningful information about the software’s trustworthiness and security.

## Developer perspective
The system is easy for developers to extend or improve, as it exposes a straightforward and well-documented REST API.

### Stack
#### Backend
- Golang
- MongoDB
- Google Cloud (Generative AI API)
- OpenAPI specification
- Swagger UI

#### Frontend
- React
- Next.js

#### Deployment
- Docker & Docker Compose
- Caddy

## Google Cloud & Gemini API Challenge

I used Google Gemini to:
- Identify the software product the user is interested in
- Search for security-related information and generate comprehensive software trust reports


## GoDaddy Challenge

I came up with a creative domain name and registered it for free using GoDaddy: softthat.fit (“soft that fits”)
