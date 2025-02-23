CREATE TABLE tags
(
    id serial PRIMARY KEY,
    name VARCHAR(255) UNIQUE
);

CREATE TABLE category
(
    id SERIAL PRIMARY KEY,
    name VARCHAR(55)
);

CREATE TABLE orders
(
    id serial PRIMARY KEY,
    petID int,
    quantity int,
    shipDate TIMESTAMP,
    status VARCHAR(10) CHECK(status IN ('placed','approved','delivered')) DEFAULT 'placed',
    complete boolean
);

CREATE TABLE pet
(
    id serial PRIMARY KEY, 
    category int,
    name VARCHAR(55),
    photoUrls VARCHAR(255),
    status VARCHAR(10) CHECK(status IN ('available','pending','sold')),
    FOREIGN KEY(category) REFERENCES category(id)
);

CREATE TABLE users
(
    id serial PRIMARY KEY,
    username varchar(20) NOT NULL UNIQUE,
    firstname varchar(20),
    lastname varchar(20),
    phone varchar(55),
    email varchar(25) NOT NULL UNIQUE,
    status int, 
    password varchar(255) NOT NULL,
    deleted_at TIMESTAMP NULL DEFAULT NULL
);

CREATE TABLE pet_tags
(
    pet_id int,
    tag_id int,
    FOREIGN KEY(pet_id) REFERENCES pet(id) ON DELETE CASCADE,
    FOREIGN KEY(tag_id) REFERENCES tags(id) ON DELETE CASCADE,
    PRIMARY KEY(pet_id, tag_id)
);


