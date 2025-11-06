-- +goose Up
CREATE TABLE pet (
    id INTEGER GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    name TEXT NOT NULL,
    dateOfBirth DATE,
    dateOfBirthExact BOOLEAN,
    imageUrl TEXT,
    about_text VARCHAR(1000),
    species TEXT,
    breed TEXT,
    sex TEXT,
    isPubliclyViewable BOOLEAN NOT NULL DEFAULT FALSE,
    likesHidden BOOLEAN NOT NULL DEFAULT FALSE,
    skillsHidden BOOLEAN NOT NULL DEFAULT FALSE,
    goalsHidden BOOLEAN NOT NULL DEFAULT FALSE,
    titlesHidden BOOLEAN NOT NULL DEFAULT FALSE,
    created_at DATE NOT NULL,
    updated_at DATE NOT NULL
);

CREATE TABLE UserPets (
    userId INTEGER NOT NULL,
    petId INTEGER NOT NULL,
    permissions_level INTEGER NOT NULL,
    active BOOLEAN NOT NULL DEFAULT FALSE,
    hidden BOOLEAN NOT NULL DEFAULT FALSE,
    created_at DATE NOT NULL,
    updated_at DATE NOT NULL,
    CONSTRAINT fk_user_pets_users
    FOREIGN KEY (userId)
    REFERENCES users(id)
    ON DELETE CASCADE,
    CONSTRAINT fk_user_pets_pets
    FOREIGN KEY (petId)
    REFERENCES pet(id)
    ON DELETE CASCADE
);

-- +goose Down
DROP TABLE UserPets;
DROP TABLE pet;