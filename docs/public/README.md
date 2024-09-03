# API Documentation

## Overview
This API offers endpoints for user retrieving movie information, and booking seats. The routes are structured to accommodate various user roles, such as admin and general users, ensuring secure access and efficient management of movie-related data.


---
### 1. Movies
- **Method**: `GET`
- **Endpoint**: `/api/v1/movies?status={status}&page={page}&size={size}`
- **Description**: Retrieves a paginated list of movies based on their status ("upcoming" or "incinemas"). The pagination parameters (page and size) are optional, with defaults applied if they are not provided.
- **Request Parameter**:
  - **Query Parameter**: 
    - status (required): Filter movies by their availability status.
        - upcoming: Movies that are not yet available for ticket booking.
        - incinemas: Movies that are currently showing and available for ticket booking.
    - page (optional): The page number of results to retrieve. Default is `1`
    - size (optional): The number of results per page. Default is `5`
- **Response Parameter**:
    - Status Code: `200 Status OK`
        ```json
        [
            {
                "id": 2,
                "title": "Batman Begins",
                "duration": 140,
                "genre": "Action, Crime, Drama",
                "release_date": "2005-07-17",
                "portrait_url": "http://localhost:8080/public/static/images/4b092dd6aed34bf0a4c64b1b5eccb242.webp",
                "landscape_url": null
            }
            //,...many more movies
        ]
   - Status Code: `400 Bad Request`
        query parameter missing or invalid
   - Status Code: `404 Not Found`
        no movies to show        
   - Status Code: `500 Internal Server Error`
        internal server error
---
### 2. Movie detials
- **Method**: `GET`
- **Endpoint**: `/api/v1/movie/{movie_id}`
- **Description**: Retrives in detials of the movie using there movie ID
- **Request Parameter**:
  - **Path Parameter**: 
    - movie_id (required): Movied ID of the movie you wish retrives the detials. `eg /api/v1/movie/1`
- **Response Parameter**:
    - Status Code: `200 Status OK`
        ```json
        {
            "id": 2,
            "title": "Batman Begins",
            "description": "Driven by tragedy, billionaire Bruce Wayne dedicates his life to uncovering and defeating the corruption that plagues his home, Gotham City. Unable to work within the system, he instead creates a new identity, a symbol of fear for the criminal underworld - The Batman.",
            "duration": 140,
            "genre": null,
            "release_date": "2005-07-17",
            "portrait_url": "http://localhost:8080/public/static/images/4b092dd6aed34bf0a4c64b1b5eccb242.webp",
            "landscape_url": "http://localhost:8080/public/static/images/f42144d5f12c4e89a718663ab2525df0.webp",
            "actors": [
                {
                    "id": 1,
                    "name": "Christian Bale",
                    "alias": "Bruce Wayne",
                    "poster": "http://localhost:8080/public/static/images/f69754c140a04054b824d90e6ef82992.webp"
                }
            ],
            "actress": [
                {
                    "id": 2,
                    "name": "Katie Holmes",
                    "alias": "Rachel Dawes",
                    "poster": "http://localhost:8080/public/static/images/b1de12e91fd4499aa03abdb17b337c90.webp"
                }
            ],
            "directors": [
                {
                    "id": 2,
                    "name": "Christopher Nolan",
                    "alias": null,
                    "poster": "http://localhost:8080/public/static/images/b209e7e677a84ea3a7b0a4c2506f69b5.webp"
                }
            ],
            "producers": [
                {
                    "id": 2,
                    "name": "Emma Thomas",
                    "alias": null,
                    "poster": "http://localhost:8080/public/static/images/2f2bae02b3744f928d65c39cbc7cf896.web"
                }
            ]
        }
   - Status Code: `400 Bad Request`
        ID path parameter missing
   - Status Code: `404 Not Found`
        if there is no movie bind to thet movie ID       
   - Status Code: `500 Internal Server Error`
        internal server error
---
### 3. Avilable Hall of that movie
- **Method**: `GET`
- **Endpoint**: `/api/v1/halls?movie_id="movie_id"`
- **Description**: Retrives the list  of avilable halls that show that movie in cinemas
- **Request Parameter**:
  - **Query Parameter**: 
    - movie_id (required): Movied ID of the movie you wish to retrives the list of avilable halls. `eg /api/v1/halls?movie_id=1`
- **Response Parameter**:
    - Status Code: `200 Status OK`
        ```json
            [
                {
                    "hall_id": "1",
                    "hall_name": "Tanthapolis Cinema",
                    "address": "Lamphelpat",
                    "city": "Imphal",
                    "state": "Manipur",
                    "postal_code": "795004",
                    "latitude": 0,
                    "longitude": -74.006
                }//, many more halls
            ]
   - Status Code: `400 Bad Request`
        query parameter missing or invalid      
   - Status Code: `500 Internal Server Error`
        internal server error
---
### 4. Avilable Timings of that movie of a hall
- **Method**: `GET`
- **Endpoint**: `/api/v1/hall/{hall_id}/showtimes?movie_id={movie_id}`
- **Description**: Retrives the list of avilable show/movie timings of that hall that you can book a ticket for it 
- **Request Parameter**:
  - **Path Parameter**: 
    - hall_id (required): hall ID of the hall you wish to retrives the list of avilable timings. 
  - **Query Parameter**: 
    - movie_id (required): Movied ID of the movie you wish to retrives the list of avilable timings. `eg /api/v1/hall/1/showtimes?movie_id=1`
- **Response Parameter**:
    - Status Code: `200 Status OK`
        ```json
        [
            {
                "date_id": 3,
                "show_date": "2024-09-10",
                "timings": [
                    {
                        "timing_id": 5,
                        "timing": "09:00:00"
                    }
                ]
            }
            //, and many more timings
        ]
   - Status Code: `204 No Content`
        if no timings are avilable for that show
   - Status Code: `400 Bad Request`
        path or query parameter missing or invalid      
   - Status Code: `500 Internal Server Error`
        internal server error
---
### 5. Hall Seat Layout
- **Method**: `GET`
- **Endpoint**: `/api/v1/hall/{hall_id}/seatlayout?timing_id=id`
- **Description**: Retrives the seat layout of that hall of that show timing 
- **Request Parameter**:
  - **Path Parameter**: 
    - hall_id (required): hall ID of the hall you wish to retrives the hall seat Layout. 
  - **Query Parameter**: 
    - timing_id (required): timing ID of the movie you wish to retrives the hall seat layout. `eg /api/v1/hall/1/seatlayout?timing_id=5`
- **Response Parameter**:
    - Status Code: `200 Status OK`
        ```json
        {
            "seat_types": [
                {
                    "id": 1,
                    "name": "Diamond",
                    "price": 350,
                    "seat_row": 4,
                    "seat_column": 8,
                    "seat_matrix": "11111111111111111111111111000011",
                    "order_from_screen": 0,
                    "row_names": [
                        "D",
                        "C",
                        "B",
                        "A"
                    ]
                },
                {
                    "id": 2,
                    "name": "Gold",
                    "price": 270,
                    "seat_row": 8,
                    "seat_column": 15,
                    "seat_matrix": "111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111",
                    "order_from_screen": 1,
                    "row_names": [
                        "H",
                        "G",
                        "F",
                        "E",
                        "D",
                        "C",
                        "B",
                        "A"
                    ]
                },
                {
                    "id": 3,
                    "name": "Silver",
                    "price": 200,
                    "seat_row": 4,
                    "seat_column": 15,
                    "seat_matrix": "111111111111111111111111111111111111111111111111111111111111111111111111111",
                    "order_from_screen": 2,
                    "row_names": [
                        "D",
                        "C",
                        "B",
                        "A"
                    ]
                }
            ],
            "reserved_seats": [
                "@"
            ],
            "booked_seats": [
                "@"
            ]
        }
   - Status Code: `400 Bad Request`
        path or query parameter missing or invalid      
   - Status Code: `500 Internal Server Error`
        internal server error
---
### 6. Checkout Hall Seats for Booking a Ticket(s) (USER)
- **Method**: `POST`
- **Endpoint**: `/api/v1/seats/checkout/{timing_id}`
- **Description**: Allows a user to check out one or multiple seats for booking. All selected seats must be of the same type. For example:
  - ✅: Seats: diamond (D9, D3, ...)
  - ❌: Seats: diamond (D9), silver (C2)
- **Note**: Users can only checkout or book seats before the pre_expiry_secs. For instance, if pre_expiry_secs = 900, users must checkout/book seats at least 15 minutes before the show starts.
- **Request Parameter**:
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
---
### 7. Book Tickets for Reserved Seats (USER)
- **Method**: `POST`
- **Endpoint**: `/api/v1/seats/book/{timing_id}`
- **Description**: Allows a user to book one or more seats that have been reserved.
- **Note**: Users can only checkout or book seats before the pre_expiry_secs. For instance, if pre_expiry_secs = 900, users must checkout/book seats at least 15 minutes before the show starts.
- **Request Parameter**:
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