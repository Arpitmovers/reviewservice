CREATE TABLE zuzu_db.Hotels (
    hotel_id BIGINT UNSIGNED  PRIMARY KEY,
    hotel_name VARCHAR(500) NOT NULL,
    platform VARCHAR(100) NOT NULL
);



CREATE TABLE zuzu_db.Providers (
    provider_id INT PRIMARY KEY,
    provider_name VARCHAR(100) NOT NULL
);

CREATE TABLE zuzu_db.Reviewers (
    reviewer_id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    display_name VARCHAR(100),
    country_id INT,
    country_name VARCHAR(100),
    review_group_id INT,
    review_group_name VARCHAR(100),
    reviewer_reviewed_count INT,
    is_expert_reviewer BOOLEAN
) ;


CREATE TABLE zuzu_db.Reviews (
    review_id BIGINT PRIMARY KEY,
    hotel_id BIGINT unsigned NOT NULL,
    provider_id INT NOT NULL,
    reviewer_id BIGINT unsigned NOT NULL,
    rating DECIMAL(3,1),
    check_in_month_year VARCHAR(20),
    review_date TIMESTAMP,
    review_title VARCHAR(200),
    review_text TEXT,
    response_text TEXT,
    room_type VARCHAR(100),
    length_of_stay INT,
    positives TEXT,
    negatives TEXT,
    encrypted_data TEXT,
    FOREIGN KEY (hotel_id) REFERENCES Hotels(hotel_id),
        FOREIGN KEY (provider_id) REFERENCES Providers(provider_id),
    FOREIGN KEY (reviewer_id) REFERENCES Reviewers(reviewer_id)
);


CREATE TABLE zuzu_db.ProviderSummary (
    hotel_id BIGINT  unsigned  NOT NULL,
    provider_id INT NOT NULL,
    overall_score DECIMAL(3,1),
    review_count INT,
    cleanliness DECIMAL(3,1),
    facilities DECIMAL(3,1),
    location DECIMAL(3,1),
    service DECIMAL(3,1),
    value_for_money DECIMAL(3,1),
    room_comfort DECIMAL(3,1),
    PRIMARY KEY (hotel_id, provider_id),
    FOREIGN KEY (hotel_id) REFERENCES Hotels(hotel_id),
    FOREIGN KEY (provider_id) REFERENCES Providers(provider_id)
);