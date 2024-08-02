CREATE TABLE profile (
    id INT AUTO_INCREMENT PRIMARY KEY,
    name VARCHAR(255),
    poster_url VARCHAR(255)
);

CREATE TABLE admin (
    id INT AUTO_INCREMENT PRIMARY KEY,
    email VARCHAR(255) UNIQUE NOT NULL,
    email_verified BOOLEAN,
    hall_registered BOOLEAN,
    profile_id INT UNIQUE,
    FOREIGN KEY (profile_id) REFERENCES profile(id)
);

CREATE TABLE user (
    id INT AUTO_INCREMENT PRIMARY KEY,
    email VARCHAR(255) UNIQUE NOT NULL,
    email_verified BOOLEAN,
    profile_id INT UNIQUE,
    FOREIGN KEY (profile_id) REFERENCES profile(id)
);

CREATE TABLE hall (
    id INT AUTO_INCREMENT PRIMARY KEY,
    hall_name VARCHAR(255) UNIQUE NOT NULL,
    hall_manager VARCHAR(255) NOT NULL,
    hall_contact VARCHAR(255) UNIQUE NOT NULL,
    admin_id INT,
    FOREIGN KEY (admin_id) REFERENCES admin(id)
);

CREATE TABLE hall_location (
    id INT AUTO_INCREMENT PRIMARY KEY,
    address VARCHAR(255) NOT NULL,
    city VARCHAR(255) NOT NULL,
    state VARCHAR(255) NOT NULL,
    postal_code VARCHAR(255) NOT NULL,
    latitude DECIMAL(9,6),
    longitude DECIMAL(9,6),
    hall_id INT,
    FOREIGN KEY (hall_id) REFERENCES hall(id)
);

CREATE TABLE hall_seat_layout (
    id INT AUTO_INCREMENT PRIMARY KEY,
    max_capacity INT NOT NULL,
    h_rows INT NOT NULL,
    h_columns INT NOT NULL,
    types VARCHAR(255),
    layout TEXT NOT NULL,
    hall_id INT,
    FOREIGN KEY (hall_id) REFERENCES hall(id)
);

CREATE TABLE seat_type (
    id INT AUTO_INCREMENT PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    start INT,
    end INT,
    price INT,
    hall_seat_layout_id INT,
    FOREIGN KEY (hall_seat_layout_id) REFERENCES hall_seat_layout(id)
);

CREATE TABLE hall_operation_time (
    id INT AUTO_INCREMENT PRIMARY KEY,
    open_time TIME NOT NULL,
    closed_time TIME NOT NULL,
    hall_id INT,
    FOREIGN KEY (hall_id) REFERENCES hall(id)
);

CREATE TABLE poster_urls (
    id INT AUTO_INCREMENT PRIMARY KEY,
    url VARCHAR(2024) NOT NULL
);

CREATE TABLE movie (
    id INT AUTO_INCREMENT PRIMARY KEY,
    title VARCHAR(255) NOT NULL,
    description TEXT NOT NULL,
    duration INT NOT NULL,
    genre VARCHAR(255),
    release_date DATE NOT NULL,
    portrait_poster_url_id INT,
    landscape_poster_url_id INT,
    FOREIGN KEY (portrait_poster_url_id) REFERENCES poster_urls(id),
    FOREIGN KEY (landscape_poster_url_id) REFERENCES poster_urls(id)
);

CREATE TABLE actor (
    id INT AUTO_INCREMENT PRIMARY KEY,
    name VARCHAR(255),
    poster_url_id INT NULL,
    FOREIGN KEY (poster_url_id) REFERENCES poster_urls(id)
);

CREATE TABLE actress (
    id INT AUTO_INCREMENT PRIMARY KEY,
    name VARCHAR(255),
    poster_url_id INT NULL,
    FOREIGN KEY (poster_url_id) REFERENCES poster_urls(id)
);

CREATE TABLE director (
    id INT AUTO_INCREMENT PRIMARY KEY,
    name VARCHAR(255),
    poster_url_id INT NULL,
    FOREIGN KEY (poster_url_id) REFERENCES poster_urls(id)
);

CREATE TABLE producer (
    id INT AUTO_INCREMENT PRIMARY KEY,
    name VARCHAR(255),
    poster_url_id INT NULL,
    FOREIGN KEY (poster_url_id) REFERENCES poster_urls(id)
);

CREATE TABLE movie_actor (
    id INT AUTO_INCREMENT PRIMARY KEY,
    movie_id INT,
    actor_id INT,
    alias VARCHAR(255),
    FOREIGN KEY (movie_id) REFERENCES movie(id),
    FOREIGN KEY (actor_id) REFERENCES actor(id)
);

CREATE TABLE movie_actress (
    id INT AUTO_INCREMENT PRIMARY KEY,
    movie_id INT,
    actress_id INT,
    alias VARCHAR(255),
    FOREIGN KEY (movie_id) REFERENCES movie(id),
    FOREIGN KEY (actress_id) REFERENCES actress(id)
);

CREATE TABLE movie_director (
    id INT AUTO_INCREMENT PRIMARY KEY,
    movie_id INT,
    director_id INT,
    FOREIGN KEY (movie_id) REFERENCES movie(id),
    FOREIGN KEY (director_id) REFERENCES director(id)
);

CREATE TABLE movie_producer (
    id INT AUTO_INCREMENT PRIMARY KEY,
    movie_id INT,
    producer_id INT,
    FOREIGN KEY (movie_id) REFERENCES movie(id),
    FOREIGN KEY (producer_id) REFERENCES producer(id)
);

CREATE TABLE movie_show (
    id INT AUTO_INCREMENT PRIMARY KEY,
    movie_id INT,
    hall_id INT,
    status VARCHAR(255) NOT NULL,
    FOREIGN KEY (movie_id) REFERENCES movie(id),
    FOREIGN KEY (hall_id) REFERENCES hall(id)
);

CREATE TABLE movie_show_dates (
    id INT AUTO_INCREMENT PRIMARY KEY,
    show_date DATE NOT NULL,
    movie_show_id INT,
    FOREIGN KEY (movie_show_id) REFERENCES movie_show(id)
);

CREATE TABLE movie_show_timings (
    id INT AUTO_INCREMENT PRIMARY KEY,
    show_timing TIME NOT NULL,
    ticket_status BOOLEAN,
    movie_show_dates_id INT NOT NULL,
    FOREIGN KEY (movie_show_dates_id) REFERENCES movie_show_dates(id)
);

CREATE TABLE booking (
    id INT AUTO_INCREMENT PRIMARY KEY,
    user_id INT,
    movie_show_id INT,
    seat_number VARCHAR(255),
    booking_timing DATETIME,
    FOREIGN KEY (user_id) REFERENCES user(id),
    FOREIGN KEY (movie_show_id) REFERENCES movie_show(id)
);

CREATE TABLE ticket (
    id INT AUTO_INCREMENT PRIMARY KEY,
    ticket_number VARCHAR(255) UNIQUE NOT NULL,
    booking_id INT,
    FOREIGN KEY (booking_id) REFERENCES booking(id)
);
