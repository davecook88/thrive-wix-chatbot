# Gemini Task: Wix Services Caching

This document outlines the changes made to the Thrive Wix Chatbot to implement caching for Wix services.

## Problem

The chatbot was fetching the list of available services from Wix every time a user initiated a conversation. This was causing a delay in the chatbot's response time, as the application had to wait for the Wix API to respond.

## Solution

To improve the chatbot's performance, I implemented a caching mechanism for the Wix services. The solution involves the following changes:

### 1. Firebase Cache

I created a new collection in Firebase called `wix-services` to store a cached version of the Wix services. This allows the application to quickly retrieve the services without having to make a call to the Wix API every time.

### 2. Background Fetching

When a user starts a conversation, the application now first checks if a cached version of the services exists in Firebase. 

- If a cached version exists and is less than 24 hours old, it is returned to the user immediately.
- If the cached version is older than 24 hours, it is returned to the user, but a background process is initiated to fetch the latest services from Wix and update the cache.
- If no cached version exists, the application returns `nil` to the user and initiates a background process to fetch the services from Wix and populate the cache.

This ensures that the user's conversation is not blocked while the services are being fetched.

### 3. Code Refactoring

I refactored the code to support the new caching mechanism. This included:

- Creating a new `db` package to handle all interactions with Firebase.
- Modifying the `main.go` file to create a single `db.Client` instance that is shared across all handlers.
- Updating the `handlers` package to use the `db.Client` to get the Wix services from the cache.

## Outcome

These changes have significantly improved the performance of the chatbot by reducing the latency of the initial response. The application is now more scalable and resilient, as it is no longer dependent on the availability of the Wix API for its core functionality.


## WIX

It is important that the chatbot understand how the classes work in WIX:

CLASS - Group class that can be booked on an ad-hoc basis or as part of a package. These tend to be conversational classes which are independent.

COURSE - Group classes that follow a specific course. Students might book an 8 session course on Wednesday at 8pm. The idea is that the group is static and all students attend all sessions.

APPOINTMENT - Ad-hoc classes that can be reserved one-on-one with a teacher. If a teacher is available, any student can book a session with them. These can be paid for individually or a package can be bought to get a discount.