# API Documentation

## Overview
This API provides endpoints for managing movie hall registrations, seat layouts, movie shows, and ticket bookings for admin specific. It includes administrative endpoints for hall management, movie show registration, and ticketing functionalities. The following sections detail each endpoint's purpose and usage.

---
## Administrative Endpoints
### 1. Register a Hall
- **Method**: `POST`
- **Endpoint**: `/api/v1/auth/admin/hall/register`
- **Description**:  Allows an admin to register a new hall with detailed information.
- **Request Parameters**:
  - Headers:
    - `Content-Type`: `application/json`
    - `WithCredentials`: `true` 
  - Body Payload: 
    ```json
    {
        "name": "Tanthapolis Cinema",
        "manager": "Foo Boo",
        "contact": "8390208210",
        "location": {
            "address": "Lamphelpat",
            "city": "Imphal",
            "state": "Manipur",
            "postal_code": "795004",
            "latitude": 40.7128, 
            "longitude": -74.0060
        },
        "operationtime": {
            "open_time": "09:00:00", //HH:MM:SS
            "close_time": "19:00:00" //HH:MM:SS
        }
    }
- **Response**:
  - Status Code: `201`
    ```json
    {
        "message": "hall registered successfully"
    }
  - Status Code: `400`
    ```json
    {
        "message": "some specific message"
    }
  - Status Code: `500`
    ```json
    {
        "message": "internal server error"
    }
  - Status Code: `409`
    ```json
    {
        "message": "a hall with that name already exists"
    }

### 2. Register a Hall SeatLayout
- **Method**: `POST`
- **Endpoint**: `/api/v1/auth/admin/hall/seatlayout/register`
- **Data Exchange Format**: `JSON`
- **Description**:  Allows an admin to register a seat layout for a hall.
- **Request Parameters**:
  - Headers:
    - `Content-Type`: `application/json`
    - `WithCredentials`: `true` 
  - Body Payload: 
    ```json
    {
        "seat_types": [
            {
                "name": "Diamond",
                "price": 350,
                "seat_row": 4,
                "seat_column": 8,
                "seat_matrix": "11111111111111111111111111000011",
                "order_from_screen": 0,
                "row_names": ["D", "C", "B", "A"]
            },
            {
                "name": "Gold",
                "price": 270,
                "seat_row": 8,
                "seat_column": 15,
                "seat_matrix": "111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111",
                "order_from_screen": 1,
                "row_names": ["H", "G", "F", "E", "D", "C", "B", "A"]
            },
            {
                "name": "Silver",
                "price": 200,
                "seat_row": 4,
                "seat_column": 15,
                "seat_matrix": "111111111111111111111111111111111111111111111111111111111111111111111111111",
                "order_from_screen": 2,
                "row_names": ["D", "C", "B", "A"]
            }
        ]
    }
- **Response**:
  - Status Code: `201`
    ```json
    {
        "message": "hall seatlayout registered successfully"
    }
  - Status Code: `400`
    ```json
    {
        "message": "some specific message"
    }
  - Status Code: `404`
    ```json
    {
        "message": "internal server error"
    }
---
### 3. Register a Movie Show
 - **Method**: `POST`
 - **Endpoint**: `/api/v1/auth/admin/hall/show/register`
 - **Description**: Allows an admin to register a movie show, including movie details, cast, and show timings.
 - **Flow**
    - **Movie Existence Check**: Check if the movie exists using the [Search Movie route](https://github.com/ShuaibKhan786/movie-ticketing-api/blob/main/docs/utility/README.md#1-search-movie-common). If it exists, use the `movie_id`. If not, register the movie.
    - **Cast Existence Check**: Check if the cast exists using the [Search Cast route](https://github.com/ShuaibKhan786/movie-ticketing-api/blob/main/docs/utility/README.md#2-search-cast-admin-only). If they exist, provide an alias; otherwise, register the cast details.
    - **Timing Conflict Check**: Verify show timings do not conflict using the [Check Timings Availability route](https://github.com/ShuaibKhan786/movie-ticketing-api/blob/main/docs/utility/README.md#3-check-timings-avilability-admin-only).
    - **Movie Poster Upload**: Prompt the admin to upload the movie poster or cast images using the [Image Upload route](https://github.com/ShuaibKhan786/movie-ticketing-api/blob/main/docs/utility/README.md#4-image-upload-admin-only).
    
  - **Request Parameter**
    - Headers:
      - `Content-Type`: `application/json`
      - `WithCredentials`: `true`
    - Body Payload:
    ```json
      {
        "status": true, // if this is true one of the timing must also be true
        "movie": {
          //"id": , //If the movie already exists, provide its ID
          "title": "Batman Begins",
          "description": "Driven by tragedy, billionaire Bruce Wayne dedicates his life to uncovering and defeating the corruption that plagues his home, Gotham City. Unable to work within the system, he instead creates a new identity, a symbol of fear for the criminal underworld - The Batman.",
          "duration": 140,//in minutes
          "genre": "Action, Crime, Drama",
          "release_date": "2005-07-17",//YYYY-MM-DD
          "portrait_url": "http://localhost:8080/public/static/images/4b092dd6aed34bf0a4c64b1b5eccb242.webp",
          "landscape_url": "http://localhost:8080/public/static/images/f42144d5f12c4e89a718663ab2525df0.webp"
        },
        "cast": {
          "actors": [
            {
              "id":  1,
              //"name": ,
              "alias": "Bruce Wayne"
              //"poster": ,
            }
          ],
          "actress": [
            {
              //"id": ,
              "name": "Katie Holmes",
              "alias": "Rachel Dawes",
              "poster": "http://localhost:8080/public/static/images/b1de12e91fd4499aa03abdb17b337c90.webp"
            }
          ],
          "directors": [
            {
              //"id": ,
              "name": "Christopher Nolan",
              "poster": "http://localhost:8080/public/static/images/b209e7e677a84ea3a7b0a4c2506f69b5.webp"
            }
          ],
          "producers": [
            {
              //"id": ,
              "name": "Emma Thomas",
              "poster": "http://localhost:8080/public/static/images/2f2bae02b3744f928d65c39cbc7cf896.web"
            }
          ]
        },
      "movie_show_timing": [
        {
          "show_date": "2024-09-10",//YYYY-MM-DD
          "show_timing": [
            {
              "time": "09:00:00",//HH:MM:SS
              "ticket_status": true,
              "pre_expiry_secs": 900,//in seconds
              "post_expiry_secs": 900//in seconds
            },
            {
              "time": "14:00:00",//HH:MM:SS
              "ticket_status": false,
              "pre_expiry_secs": 900,//in seconds
              "post_expiry_secs": 900//in seconds
            }
          ]
        },
        {
          "show_date": "2024-09-11",//YYYY-MM-DD
          "show_timing": [
            {
              "time": "09:00:00",//HH:MM:SS
              "ticket_status": false,
              "pre_expiry_secs": 900,//in seconds
              "post_expiry_secs": 900//in seconds
            },
            {
              "time": "14:00:00",//HH:MM:SS
              "ticket_status": false,
              "pre_expiry_secs": 900,//in seconds
              "post_expiry_secs": 900//in seconds
            }
          ]
        }
      ]
    } 



  - **Response**
   - Status Code: `201 Created`
        successfully register a show
   - Status Code: `400 Bad Request`
        invalid body payload / hall not registered / hall seat layout not registered
   - Status Code: `401 Unauthorized`
        inavlid access token or expired (so renew it)
   - Status Code: `500 Internal Server Error`
        internal server error

### 4. Get All Registered Movies/Shows of an Admin
 - **Method**: `GET`
 - **Endpoint**: `/api/v1/auth/admin/hall/shows`
 - **Description**: Allows an admin to retrieves all registered movies or shows for an admin.
 - **Request Parameter**
    - Headers:
      - `Content-Type`: `application/json`
      - `WithCredentials`: `true`
  - **Response**
   - Status Code: `200 Status OK`
      ```json
      [
        {
            "show_id": 1,
            "movie_id": 1,
            "status": true,
            "movie_title": "Batman Begins",
            "movie_poster_url": "http://localhost:8080/public/static/images/4b092dd6aed34bf0a4c64b1b5eccb242.webp"
        }
        //,and many more movies
      ]
   - Status Code: `400 Bad Request`
      hall not registered / hall seat layout not registered
   - Status Code: `401 Unauthorized`
      inavlid access token or expired (so renew it)
   - Status Code: `500 Internal Server Error`
      internal server error
### 5. Get All Timings for a Movie/Show
 - **Method**: `GET`
 - **Endpoint**: `/api/v1/auth/admin/hall/show/{show_id}/timings`
 - **Description**: Allows an admin to retrieved all registered timings for a specific movie or show.
 - **Request Parameter**
    - **Headers**:
      - `Content-Type`: `application/json`
      - `WithCredentials`: `true`
    - **Path Parameter**:
      - show_id (required): The ID of the show for which you want to retrieve timings. `eg: /api/v1/auth/admin/hall/show/1/timings`
  - **Response**
   - Status Code: `200 Status OK`
      ```json
      [
        {
            "show_date": "2024-09-10",
            "timings": [
                {
                    "timing_id": 1,
                    "timing": "09:00:00",
                    "pre_expiry_secs": 900,
                    "post_expiry_secs": 900,
                    "ticket_status": true
                },
                {
                    "timing_id": 2,
                    "timing": "14:00:00",
                    "pre_expiry_secs": 900,
                    "post_expiry_secs": 900,
                    "ticket_status": false
                }
            ]
        },
        {
            "show_date": "2024-09-11",
            "timings": [
                {
                    "timing_id": 3,
                    "timing": "09:00:00",
                    "pre_expiry_secs": 900,
                    "post_expiry_secs": 900,
                    "ticket_status": false
                },
                {
                    "timing_id": 4,
                    "timing": "14:00:00",
                    "pre_expiry_secs": 900,
                    "post_expiry_secs": 900,
                    "ticket_status": false
                }
            ]
        }
        //,many more timings
      ]
   - Status Code: `400 Bad Request`
      missing / invalid path parameter
   - Status Code: `401 Unauthorized`
      inavlid access token or expired (so renew it)
   - Status Code: `500 Internal Server Error`
      internal server error
### 6. Release a Show Timing for Ticket Booking
 - **Method**: `POST`
 - **Endpoint**: `/api/v1/auth/admin/hall/show/ticket/release/{timing_id}`
 - **Description**: This endpoint allows an admin to release a specific show timing, making it available for ticket booking by both users and the admin.
 - **Request Parameter**
   - **Headers**:
      - `Content-Type`: `application/json`
      - `WithCredentials`: `true`
    - **Path Parameter**:
      - timing_id (required): The ID of the show timing you wish to release for ticket booking. `eg: /api/v1/auth/admin/hall/show/ticket/release/1`
  - **Response**
   - Status Code: `200 Status OK`
      successfully released a specific timing for ticket booking
   - Status Code: `400 Bad Request`
      missing / invalid path parameter
   - Status Code: `401 Unauthorized`
      inavlid access token or expired (so renew it)
   - Status Code: `500 Internal Server Error`
      internal server error

---
## Book Tickets for Client (Admin)
### 1. Checkout Hall Seats for Booking a Ticket(s) (ADMIN)
- **Method**: `POST`
- **Endpoint**: `/api/v1/auth/admin/seats/checkout/{timing_id}`
- **Description**:  Allows an admin to check out one or multiple seats for booking on behalf of a client. All selected seats must be of the same type. For example
  - ✅: Seats: diamond (D9, D3, ...)
  - ❌: Seats: diamond (D9), silver (C2)
- **Note**: Admins have the ability to book seats for users even after the show has started, up until the show's start time plus the post_expiry_secs. For example, if `post_expiry_secs = 900`, admins can book seats up to 15 minutes after the show has begun.
- **Request Parameter**:
  - **Headers**:
    - `Content-Type`: `application/json`
    - `WithCredentials`: `true`
  - **Path Prameter**
    - timing_id (required): timing ID of the show you wish to reserved seats.
  - **Body Payload**
    ```json
    {
        "id": 1, // seat_type ID
        "counts": 2, // number of seats to be booked
        "seats": ["D9", "C2"] // list of seat identifiers
    }

- **Response Parameter**:
    - Status Code: `200 Status OK`
        ```json
        {
            "id": 1,
            "counts": 2,
            "seats": [
                "D9",
                "C2"
            ],
            "name": "Diamond",
            "price": 350,
            "movie_name": "Batman Begins",
            "hall_name": "Tanthapolis Cinema",
            "show_date": "2024-09-10",
            "show_timing": "09:00:00",
            "total_amount": 700
        }     
   - Status Code: `400 Bad Request`
         - seats are already reserved or booked by other users
            ```json
            {
                "seats": [
                    "D9",
                    "C2"
                ]
            }// here return seats are already taken by other users
        - invalid body payload      
   - Status Code: `500 Internal Server Error`
        internal server error
### 2. Book Tickets for Reserved Seats (ADMIN)
- **Method**: `POST`
- **Endpoint**: `/api/v1/auth/admin/seats/book/{timing_id}`
- **Description**: Allows an admin to book one or more seats that have been reserved for a client.
- **Note**: Admins have the ability to book seats for users even after the show has started, up until the show's start time plus the post_expiry_secs. For example, if post_expiry_secs = 900, admins can book seats up to 15 minutes after the show has begun.
- **Request Parameter**:
  - **Headers**:
    - `Content-Type`: `application/json`
    - `WithCredentials`: `true`
  - **Path Prameter**
    - timing_id (required): The timing ID of the show for which you wish to book tickets.
  - **Body Payload**
    ```json
	{
		"id": 1, // seat_type ID
		"counts": 2,
		"seats":  [
            "D9",
            "C2"
        ],                  
		"payable_amount": 700,
		"payment_mode": "upi",
		"customer_phone_no": "7301271044"
	}

- **Response Parameter**:
    - Status Code: `201 Status Created`
        ```json
        {
            "customer_phone_no": "7301271044",
            "movie_name": "Batman Begins",
            "hall_name": "Tanthapolis Cinema",
            "show_date": "2024-09-10",
            "show_time": "09:00:00",
            "tickets": [
                {
                    "ticket_number": "f00ad563-578d-40cd-90b4-01920881dc73",
                    "seat_number": "D9"
                },
                {
                    "ticket_number": "b4d78804-10f7-4ef1-b27d-914b1be5121a",
                    "seat_number": "C2"
                }
            ]
        }     
   - Status Code: `400 Bad Request`
         - If the reserved seats have expired (i.e., the TTL has expired)
            ```json
            {
                "message": "one or more seat reservations have expired"
            }
        - If the reserved seats have expired and have been booked by another user
            ```json
            {
                "message": "one or more seat are booked"
            }
        - invalid body payload      
   - Status Code: `500 Internal Server Error`
        internal server error