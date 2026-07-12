# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

### Added

- **New ACL Resources**
  - `defguard_acl_alias` - Resource for managing ACL aliases in DefGuard
  - `defguard_acl_destination` - Resource for managing ACL destinations in DefGuard
  - `defguard_acl_rule` - Resource for managing ACL rules in DefGuard

- **OpenID Provider Support**
  - `defguard_openid_provider` - Resource for managing OpenID Connect providers in DefGuard
  - Supports multiple provider kinds: Custom, Google, Microsoft, Okta, JumpCloud, Zitadel
  - Directory synchronization configuration support

- **Enhanced Authentication Methods**
  - Cookie-based authentication support via `cookie` provider field
  - Username/password authentication with automatic session management via `username` and `password` fields
  - Cookie jar support for proper session handling using `golang.org/x/net/cookiejar`
  - Session cookie name configuration (`defguard_session`)
  - Direct session value setting capability

- **Device Resource Enhancements**
  - `device_type` field (computed) - Device type (user or network)
  - `configured` field (computed) - Whether device is configured and ready to use
  - `description` field (optional+computed) - Device description

- **User Resource Enhancements**
  - `name` field (computed) - User's full name (first + last)
  - `password` field (optional, sensitive) - Password with validation requirements
    - Minimum 10 characters
    - Requires lowercase letters, uppercase letters, numbers, and special symbols

- **Group Resource Enhancements**
  - `members` now optional+computed instead of required
  - Import functionality updated to support both ID and name-based imports

### Changed

- **Provider Configuration**
  - Added validation to prevent using `api_token` and `cookie` together
  - Provider configuration now supports three authentication methods:
    - API token only
    - Cookie only  
    - Username/password (automatically acquires session)

- **Makefile Targets**
  - `deploy-dev` - Deploy to dev environment using prod folder
  - `deploy-prod` - Deploy to production environment
  - `validate` - Validate Terraform configuration in prod folder
  - `test-e2e` - Run full E2E test suite from prod folder
  - `dev-setup` - Show development environment variable setup

- **Documentation Updates**
  - README.md completely rewritten with comprehensive automation documentation
  - Updated provider resource schemas with clearer field descriptions
  - Added CIDR format specifications for network-related fields
  - Documented all new resources and authentication options

### Fixed

- **User Resource**
  - Improved error handling with more detailed API response information
  - Better handling of computed fields that may not be returned in initial create response
  - Read-back functionality to accurately populate all computed fields after creation

- **Device Resource**
  - Fixed path for device creation to use username-based endpoint (`/api/v1/device/user/{username}`)
  - Improved update logic based on OpenAPI schema

- **Group Resource**
  - Fixed POST `/api/v1/group` not returning group ID (now fetches full info after creation)
  - Updated read endpoint from `/api/v1/group-info/{name}` to `/api/v1/group/{name}`
  - Improved update logic with proper rename and member modification handling

- **Network Resource**
  - DNS field now properly handles null/unknown values using `nilIfUnknown` helper
  - Updated payload structure based on OpenAPI schema

### Technical Improvements

- Added `nilIfUnknown` helper function to handle optional string fields
- Enhanced client with cookie jar support for session management
- Added comprehensive debug logging in HTTP requests (can be removed for production)
- Improved error handling throughout all resources
- Added proper plan modifiers for computed fields