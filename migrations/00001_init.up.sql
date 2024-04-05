CREATE TABLE Cars (
                      ID           SERIAL PRIMARY KEY,
                      regNum VARCHAR(10),
                      mark VARCHAR(50),
                      model VARCHAR(50),
                      year INT,
                      owner_name VARCHAR(50),
                      owner_surname VARCHAR(50),
                      owner_patronymic VARCHAR(50)
);
