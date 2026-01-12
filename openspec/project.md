# Project Context

## Purpose

I want to build a worknote. It is an app to help workers, there are few initial core features:

- Help user log their job hunting activities, this includes the hiring journey.
- Help user log their daily activities during work.

## Tech Stack

- golang
- postgresql
- redis

## Project Conventions

### Code Style

- Use golang, minimal usage of interface
- Do not do dependency injection
- No test needed

### Architecture Patterns

- initialize database connection and prepare named during main
- a single main.go file for routing
- handler layer which will be mapped to routes
- service layer which will contain business logic
- repository layer which will contain database logic
- repository layer will have initialize() function that must be called on main
- handler layer will call service layer
- service layer will call repository layer
- use direct access, for example service layer can simply call user_repository.GetByID()
