# API Documentation

## Overview
This API provides a set of utility routes designed for both admin and user roles, catering to specific needs. 

---
### Session Management Routes

### 1. Sign In
- **Method**: `POST`
- **Endpoint**: `/api/v1/oauth/provider/signin`
- **Description**: Allows users(admin/user) to sign in using an OAuth provider.
- **Request Example**:

  ```javascript
    try {
        const response = await axios.post('http://localhost:3090/api/v1/oauth/provider/signin', {
            role: "admin",
            provider: "google",
            redirected_url: "your_redirected_url"
        }, {
            withCredentials: true, 
            headers: {
            'Content-Type': 'application/json'
            }
        });
    
        switch (response.status) {
            case 400:
                //handle bad request 
            break;
            case 202:
                const url = response.data.url;
                window.location.href = url; 
            break;
            default:
            break;
        }
    } catch (error) {
        //handle error
    }
### 2. Refreshing access token
 - **Method**: `GET`
 - **Endpoint**: `/api/v1/auth/refresh/token`
 - **Description**: To refresh the expired access JWT token if the refresh token is valid
 - **Request Parameter**:
   - Headers:
     - `WithCredentials`: `true`
 - **Response**
   - Status Code: `204`
        successfully refresh the access token
   - Status Code: `400`
        might be missing necessary cookies or access token is still valid
   - Status Code: `401`
        session expired / refresh token expired so do a relogin
   - Status Code: `403`
        access or refersh token is tampered

---
### 1. Search movie (common)
 - **Method**: `GET`
 - **Endpoint**: `/api/v1/movie?search_title="movie_name"`
 - **Description**: Allows users and admins to search for a movie. Users can check if a movie is available for booking, while admins can verify if a movie exists when registering a show.
 - **Request Parameter**:
   - **Query Parameter**: 
     - search_title(required): The movie title to search for.
          ```query
          /api/v1/movie?search_title=Batman
 - **Response** 
   - Status Code: `302 Status Found`
        ```json
        [
            {
                "id": 1,
                "status": true,
                "title": "Batman Begins",
                "portrait_url": "http://localhost:8080/public/static/images/4b092dd6aed34bf0a4c64b1b5eccb242.webp"
            },
            // ...additional movies
        ]
   - Status Code: `400 Bad Request`
        query parameter missing
   - Status Code: `404 Not Found`
        movie not found        
   - Status Code: `500 Internal Server Error`
        internal server error
---
### 2. Search Cast (admin only)
 - **Method**: `GET`
 - **Endpoint**: `/api/v1/auth/admin/cast?search_role="role"&search_name="cast_name"`
 - **Description**: Allow admin to search a cast
 - **Request Parameter**:
   - **Headers**:
      - `WithCredentials`: `true` 
   - **Query Parameter**: 
     - search_role(required): the cast role it can be `actor`, `actresss`, `director`, `producer`
     - search_name(required): the cast name you want to search e.g. <i>Christian Bale</i>
          ```query
          /api/v1/auth/admin/cast?search_role=actor&search_name=Christian Bale
 - **Response** 
   - Status Code: `302 Status Found`
        ```json
        [
            {
                "id": 1,
                "name": "Christian Bale",
                "poster_url":  "http://localhost:8080/public/static/images/f69754c140a04054b824d90e6ef82992.webp"//null
            },
            //...many more cast list
        ]
   - Status Code: `400 Bad Request`
        query parameter missing
   - Status Code: `401 Unauthorized`
        inavlid access token or expired (so renew it)
   - Status Code: `404 Not Found`
        cast not found        
   - Status Code: `500 Internal Server Error`
        internal server error
---
### 3. Check Timings Avilability (admin only)
 - **Method**: `POST`
 - **Endpoint**: `/api/v1/auth/admin/hall/show/timings/avilability`
 - **Description**: Allow admin to check wether that show timings is avilable to register so that it can prevent from timing conficts
 - **Request Parameter**:
   - **Headers**:
      - `Content-Type`: `application/json`
      - `WithCredentials`: `true` 
   - **Body Payload**: 
     ```json
     [
          {
               "show_date": "2024-07-23", //YYYY:MM:DD
               "show_timings": [
                    "09:00:00", //HH:MM:SS
                    "12:00:00" //HH:MM:SS
               ]
          }
          //...many more timings
     ]
     //dates must be send in ascending order e.g. : 2024-07-23, 2024-07-24, 2024-07-25
 - **Response** 
   - Status Code: `204 No Content`
        no timings are conflict
   - Status Code: `400 Bad Request`
        invalid body payload / hall not registered
   - Status Code: `401 Unauthorized`
        inavlid access token or expired (so renew it)
   - Status Code: `409 Conflict` the timings that are conflict with other movies 
     ```json
     [
          {
               "show_date": "2024-07-23", //YYYY:MM:DD
               "show_timings": [
                    "09:00:00" //HH:MM:SS
               ]
          }
          //...many more timings
     ]
   - Status Code: `500 Internal Server Error`
        internal server error
---
### 4. Image Upload (admin only)
 - **Method**: `POST`
 - **Endpoint**: `/api/v1/auth/admin/image/upload`
 - **Description**: Allow admin to upload an image such movie, cast poster.
 - **Request Parameter**:
   - **Headers**:
      - `Content-Type`: `multipartform/data`
      - `WithCredentials`: `true` 
   - **Body Payload**: 
     - image (file): The image file to be uploaded (required).
     - layout (string): The layout type for the image. Supported values are "portrait", "landscape", and "passport" (required).
     - format (string): The image format (e.g., "jpeg", "png", "webp")(required).
     ```plain
     image: <file>
     layout: "portrait"
     format: "webp"
 - **Response** 
   - Status Code: `201 Created`
        ```json
          {
               "image-url": "http://localhost:8080/public/static/images/f69754c140a04054b824d90e6ef82992.webp"
          }
   - Status Code: `400 Bad Request`
        invalid/missing form data
   - Status Code: `401 Unauthorized`
        inavlid access token or expired (so renew it)
   - Status Code: `500 Internal Server Error`
        internal server error

