# Project Rules and Guidelines

## Code Organization

### Directory Structure 
backend/
├── cmd/
│ └── api/ # Application entrypoint
├── internal/ # Private application code
│ ├── config/ # Configuration
│ ├── model/ # Data models
│ ├── repository/ # Data access layer
│ ├── service/ # Business logic
│ ├── handler/ # HTTP handlers
│ └── middleware/ # HTTP middleware
├── pkg/ # Public libraries
└── docs/ # Documentation



### Coding Standards

1. **Naming Conventions**
   - Use PascalCase for exported names (public)
   - Use camelCase for internal names (private)
   - Use snake_case for database fields
   - Prefix interface names with 'I' (e.g., ITranslationService)

2. **Error Handling**
   - Always check errors
   - Use custom error types for domain-specific errors
   - Log errors at the appropriate level
   - Return wrapped errors with context

3. **Comments and Documentation**
   - Write godoc comments for all exported types and functions
   - Include examples in documentation
   - Keep comments up to date with code changes

4. **Testing**
   - Write unit tests for all packages
   - Use table-driven tests where appropriate
   - Aim for >80% test coverage
   - Include integration tests for critical paths

5. **Database**
   - Use migrations for schema changes
   - Include indexes for frequently queried fields
   - Use foreign key constraints
   - Keep transactions short and focused

6. **API Design**
   - Use RESTful conventions
   - Version APIs in URL (/api/v1/...)
   - Use proper HTTP methods and status codes
   - Include validation for all inputs

7. **Security**
   - Sanitize all user inputs
   - Use prepared statements for SQL
   - Implement rate limiting
   - Use secure headers
   - Log security-related events

8. **Performance**
   - Use connection pooling
   - Implement caching where appropriate
   - Monitor and optimize database queries
   - Profile code in development

9. **Monitoring**
   - Include metrics for key operations
   - Use structured logging
   - Monitor error rates
   - Track response times

10. **Development Process**
    - Use feature branches
    - Write clear commit messages
    - Review all code changes
    - Update documentation with changes