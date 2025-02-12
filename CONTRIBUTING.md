# Contributing to goschedviz

First off, thank you for considering contributing to goschedviz! This is a great opportunity to learn about Go scheduler
internals and get experience with open source development.

_Remember: every expert was once a beginner. Don't hesitate to ask questions and make mistakes - that's how we learn!
❤️_

## Before You Start Contributing

To save time and avoid duplicate efforts, please follow these steps:

1. **Check Existing Issues and PRs**
    - Browse through open and closed issues
    - Review current Pull Requests
    - Use the search function with relevant keywords to find similar proposals

2. **Discuss Major Changes**
    - For small fixes (typos, minor bugs) you can create a PR directly
    - For substantial changes, start by creating an Issue first:
        - Describe what you want to implement
        - Explain why it's valuable
        - Outline your implementation approach
    - Wait for maintainers' feedback
    - This helps to:
        - Ensure the changes align with project goals
        - Get early implementation guidance
        - Avoid wasted effort
        - Align on technical decisions

3. **Get Confirmation**
    - Wait for maintainers' approval comment
    - Clarify any uncertainties
    - Only proceed with implementation after getting the green light

This approach helps:

- Save everyone's time
- Keep the project coherent
- Prevent conflicting changes
- Ensure higher code quality

## Getting Started

1. Fork the repository
2. Clone your fork: `git clone https://github.com/your-username/goschedviz`
3. Add upstream remote: `git remote add upstream https://github.com/JustSkiv/goschedviz`
4. Create a branch for your changes: `git checkout -b feature/your-feature-name`

## Development Environment

1. Make sure you have Go 1.23+ installed
2. Build the project: `make build`
3. Run tests: `make test`
4. Try running examples: `./bin/goschedviz -target=examples/simple/main.go`

## What Can I Contribute?

Here are some ways you can help:

### 🐛 Bug Fixes

- Run the tool with different Go programs and find edge cases
- Improve error handling
- Fix race conditions
- Make UI more stable

### 🎨 UI Improvements

- Add new visualization widgets
- Improve existing widgets layout
- Add color themes
- Make UI more responsive

### 📝 Documentation

- Improve documentation clarity
- Add more examples
- Translate documentation
- Add comments to complex parts of code

### ✨ New Features

- Add new metrics visualization
- Add export functionality
- Add configuration options
- Add profiling features

## Pull Request Process

1. Update documentation if needed
2. Add or update tests for your changes
3. Follow existing code style
4. Make sure all tests pass
5. Create a Pull Request with clear description
6. Link any related issues

## Code Style

- Follow standard Go conventions (use `gofmt`)
- Add comments for non-obvious code
- Keep functions focused and small
- Use meaningful variable names

## First Time Contributors

Never contributed to open source? Here's a simple process:

1. Find an issue labeled `good-first-issue`
2. Comment on the issue that you'd like to work on it
3. Follow the steps above to submit your changes
4. Ask questions if you're stuck - we're here to help!

## Community Guidelines

- Be respectful and constructive in discussions
- Explain your thoughts clearly
- Ask questions if something is unclear
- Help others learn

## Need Help?

- Check existing issues and discussions
- Create a new issue with questions
- Contact project maintainers via Telegram channels
- Join the community discussion
