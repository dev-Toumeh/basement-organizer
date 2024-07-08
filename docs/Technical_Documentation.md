# Technical Documentation


<!--toc:start-->
- [Technical Documentation](#technical-documentation)
  - [Decisions](#decisions)
  - [Open questions/tasks](#open-questionstasks)
  - [Overview](#overview)
  - [Technologies used](#technologies-used)
<!--toc:end-->

## Decisions
- project's main focus is learning `Go` in backend and frontend in general (starting with `HTMX` as a lightweight library)
- solve a real world problem
- other goals are simplicity, fast and simple development, desktop and mobile support
    - only browser based, no native apps
    - focus on necessary core functionalities
    - avoiding features that are nice to have
    - less frameworks/dependencies
- have more time for relevant problems and learning
    - use Github as main platform (central place for information)
        - version control
        - documentation
        - issues 
        - no infrastructure needed
        - less searching
        - ability to use CI/CD (Github Actions) at a later time if necessary

## Open questions/tasks
See [Github Issues](https://github.com/dev-Toumeh/basement-organizer/issues)

## Overview
For detailed architecture with descriptions of functions data etc. see [Detailed Architecture](detailed_architecture.md)

- Login Page
    - userame/pw
    - login process
        - requests?
        - form?
        - direct/json?
- Register user
    1. check if users.json exsits
    2. check if username exists
    3. hash password
    4. insert user in users.json and save

    - initial super admin account?
    - process?
    - db connections?
        - db schema?
    - error responses?
- Main Page (overview)
    - Create items
        - json/db ?
        - Ability to add and change descriptions to contents (maybe also images, documents like manuals, documentation, invoices, etc.).
    - create Boxes (or any kind of container) can be assigned to any location.
    - define places, like rooms.
    - Ability to see what's inside a container.

    - Ability to assign contents/items to containers?

- Ability to search for contents (e.g., searching for `keyboard` should return the box and its location).
- Capable of printing labels with a label printer (preferably not manually).
- Ability to generate QR codes or read and import existing QR codes.

## Technologies used
- Go backend
- HTMX frontend

