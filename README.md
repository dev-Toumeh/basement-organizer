# basement-organizer


<!--toc:start-->
- [basement-organizer](#basement-organizer)
  - [The Problem](#the-problem)
  - [The Idea](#the-idea)
  - [Example](#example)
  - [What it Should Do](#what-it-should-do)
  - [Goal](#goal)
  - [Start Development](#start-development)
<!--toc:end-->

Overview, detailed information and decisions can be found in [ Technical Documentation ](docs/Technical_Documentation.md).

## The Problem
- When sorting things into boxes at home and storing them in the basement, over time it's hard to remember where things are and in which box. I also don't want to constantly go to the basement to open every box searching for something.

## The Idea
- To be able to store boxes with contents in the basement and find the contents later.
- Or to virtually look through the contents of the boxes without actually finding and opening them.

## Example
- I have a box containing an old keyboard, an old router, and some hard drives.
- I enter these items into the app.
- I set a label (e.g., "old devices") and specify where the box is located (e.g., basement).
- Then, a label with a QR code can be generated for printing, which can later be scanned with a smartphone to know what's inside the box without opening it.

## What it Should Do
- Web app - mobile and desktop browser compatible.
- Should be synchronized.
- Can be self-hosted.
- Ability to generate QR codes or read and import existing QR codes.
- Capable of printing labels with a label printer (preferably not manually).
- Ability to search for contents (e.g., searching for `keyboard` should return the box and its location).
- Ability to define places, like rooms.
- Boxes (or any kind of container) can be assigned to any location.
- Ability to see what's inside a container.
- Ability to assign contents to containers.
- Ability to add and change descriptions to contents (maybe also images, documents like manuals, documentation, invoices, etc.).
- offline use (no internet acces) ? https://github.com/dev-Toumeh/basement-organizer/issues/3

## Goal
- To develop a simple way to quickly enter and find things.
- For home use.
- Not overloaded with features.
- Whether on the go with a smartphone or at home on a PC.

## Start Development
For development see [Getting Started](docs/getting_started.md)

