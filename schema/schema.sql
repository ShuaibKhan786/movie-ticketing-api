CREATE DATABASE IF NOT EXISTS movieticketing;
USE movieticketing;

CREATE TABLE `profile` (
  `id` INT PRIMARY KEY AUTO_INCREMENT,
  `name` VARCHAR(255),
  `poster_url` VARCHAR(255)
);

CREATE TABLE `admin` (
  `id` INT PRIMARY KEY AUTO_INCREMENT,
  `email` VARCHAR(255) UNIQUE NOT NULL,
  `email_verified` BOOL,
  `hall_registered` BOOL,
  `profile_id` INT UNIQUE
);

CREATE TABLE `user` (
  `id` INT PRIMARY KEY AUTO_INCREMENT,
  `email` VARCHAR(255) UNIQUE NOT NULL,
  `email_verified` BOOL,
  `profile_id` INT UNIQUE
);

CREATE TABLE `hall` (
  `id` INT PRIMARY KEY AUTO_INCREMENT,
  `hall_name` VARCHAR(255) UNIQUE NOT NULL,
  `hall_manager` VARCHAR(255) NOT NULL,
  `hall_contact` VARCHAR(255) UNIQUE NOT NULL,
  `admin_id` INT
);

CREATE TABLE `hall_location` (
  `id` INT PRIMARY KEY AUTO_INCREMENT,
  `address` VARCHAR(255) NOT NULL,
  `city` VARCHAR(255) NOT NULL,
  `state` VARCHAR(255) NOT NULL,
  `postal_code` VARCHAR(255) NOT NULL,
  `latitude` DECIMAL(9,6),
  `longitude` DECIMAL(9,6),
  `hall_id` INT
);

CREATE TABLE `hall_seat_layout` (
  `id` INT PRIMARY KEY AUTO_INCREMENT,
  `max_capacity` INT NOT NULL,
  `h_rows` INT NOT NULL,
  `h_columns` INT NOT NULL,
  `types` VARCHAR(255),
  `layout` TEXT NOT NULL,
  `hall_id` INT
);

CREATE TABLE `seat_type` (
  `id` INT PRIMARY KEY AUTO_INCREMENT,
  `name` VARCHAR(255),
  `start` INT NOT NULL,
  `end` INT NOT NULL,
  `price` INT NOT NULL,
  `hall_seat_layout` INT NOT NULL,
);

CREATE TABLE `hall_operation_time` (
  `id` INT PRIMARY KEY AUTO_INCREMENT,
  `open_time` TIME NOT NULL,
  `closed_time` TIME NOT NULL,
  `hall_id` INT
);

CREATE TABLE `movie` (
  `id` INT PRIMARY KEY AUTO_INCREMENT,
  `title` VARCHAR(255) NOT NULL,
  `description` TEXT NOT NULL,
  `duration` INT NOT NULL,
  `genre` VARCHAR(255),
  `release_date` DATE NOT NULL
);

CREATE TABLE `actor` (
  `id` INT PRIMARY KEY AUTO_INCREMENT,
  `name` VARCHAR(255)
);

CREATE TABLE `actress` (
  `id` INT PRIMARY KEY AUTO_INCREMENT,
  `name` VARCHAR(255)
);

CREATE TABLE `director` (
  `id` INT PRIMARY KEY AUTO_INCREMENT,
  `name` VARCHAR(255)
);

CREATE TABLE `producer` (
  `id` INT PRIMARY KEY AUTO_INCREMENT,
  `name` VARCHAR(255)
);

CREATE TABLE `movie_actor` (
  `id` INT PRIMARY KEY AUTO_INCREMENT,
  `movie_id` INT,
  `actor_id` INT
  `alias` VARCHAR(255) NULL
);

CREATE TABLE `movie_actress` (
  `id` INT PRIMARY KEY AUTO_INCREMENT,
  `movie_id` INT,
  `actress_id` INT
  `alias` VARCHAR(255) NULL
);

CREATE TABLE `movie_director` (
  `id` INT PRIMARY KEY AUTO_INCREMENT,
  `movie_id` INT,
  `director_id` INT
);

CREATE TABLE `movie_producer` (
  `id` INT PRIMARY KEY AUTO_INCREMENT,
  `movie_id` INT,
  `producer_id` INT
);

CREATE TABLE `movie_show` (
  `id` INT PRIMARY KEY AUTO_INCREMENT,
  `movie_id` INT,
  `hall_id` INT,
  `status` VARCHAR(255) NOT NULL,
);

CREATE TABLE `movie_show_dates` (
  `id` INT PRIMARY KEY AUTO_INCREMENT,
  `show_date` DATE NOT NULL,
  `movie_show_id` INT
);

CREATE TABLE `movie_show_timings` (
  `id` INT PRIMARY KEY AUTO_INCREMENT,
  `show_timing` TIME NOT NULL,
  `movie_show_dates_id` INT NOT NULL
);

CREATE TABLE `booking` (
  `id` INT PRIMARY KEY AUTO_INCREMENT,
  `user_id` INT,
  `movie_show_id` INT,
  `seat_number` VARCHAR(255),
  `booking_timing` DATETIME
);

CREATE TABLE `ticket` (
  `id` INT PRIMARY KEY AUTO_INCREMENT,
  `ticket_number` VARCHAR(255) UNIQUE NOT NULL,
  `booking_id` INT
);

ALTER TABLE `admin` ADD FOREIGN KEY (`profile_id`) REFERENCES `profile` (`id`);

ALTER TABLE `user` ADD FOREIGN KEY (`profile_id`) REFERENCES `profile` (`id`);

ALTER TABLE `hall` ADD FOREIGN KEY (`admin_id`) REFERENCES `admin` (`id`);

ALTER TABLE `hall_location` ADD FOREIGN KEY (`hall_id`) REFERENCES `hall` (`id`);

ALTER TABLE `hall_seat_layout` ADD FOREIGN KEY (`hall_id`) REFERENCES `hall` (`id`);

ALTER TABLE `seat_type` ADD FOREIGN KEY (`hall_seat_layout_id`) REFERENCES `hall_seat_layout` (`id`);

ALTER TABLE `hall_operation_time` ADD FOREIGN KEY (`hall_id`) REFERENCES `hall` (`id`);

ALTER TABLE `movie_actor` ADD FOREIGN KEY (`movie_id`) REFERENCES `movie` (`id`);

ALTER TABLE `movie_actor` ADD FOREIGN KEY (`actor_id`) REFERENCES `actor` (`id`);

ALTER TABLE `movie_actress` ADD FOREIGN KEY (`movie_id`) REFERENCES `movie` (`id`);

ALTER TABLE `movie_actress` ADD FOREIGN KEY (`actress_id`) REFERENCES `actress` (`id`);

ALTER TABLE `movie_director` ADD FOREIGN KEY (`movie_id`) REFERENCES `movie` (`id`);

ALTER TABLE `movie_director` ADD FOREIGN KEY (`director_id`) REFERENCES `director` (`id`);

ALTER TABLE `movie_producer` ADD FOREIGN KEY (`movie_id`) REFERENCES `movie` (`id`);

ALTER TABLE `movie_producer` ADD FOREIGN KEY (`producer_id`) REFERENCES `producer` (`id`);

ALTER TABLE `movie_show` ADD FOREIGN KEY (`movie_id`) REFERENCES `movie` (`id`);

ALTER TABLE `movie_show` ADD FOREIGN KEY (`hall_id`) REFERENCES `hall` (`id`);

ALTER TABLE `movie_show_timing` ADD FOREIGN KEY (`movie_show_id`) REFERENCES `movie_show` (`id`);

ALTER TABLE `booking` ADD FOREIGN KEY (`user_id`) REFERENCES `user` (`id`);

ALTER TABLE `booking` ADD FOREIGN KEY (`movie_show_id`) REFERENCES `movie_show` (`id`);

ALTER TABLE `ticket` ADD FOREIGN KEY (`booking_id`) REFERENCES `booking` (`id`);

ALTER TABLE `movie_producer` ADD FOREIGN KEY (`movie_id`) REFERENCES `movie_producer` (`id`);
