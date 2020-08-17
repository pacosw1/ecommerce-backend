CREATE TABLE Products (
    productID INT NOT NULL AUTO_INCREMENT,
    name VARCHAR(100) NOT NULL,
    description TEXT NOT NULL,
    stock INT NOT NULL,
    price FLOAT NOT NULL,
    salePrice FLOAT NOT NULL,
    created DATETIME NOT NULL,
    PRIMARY KEY(productID)
)


CREATE TABLE Images (
    imageID INT NOT NULL AUTO_INCREMENT,
    productID INT NOT NULL,
    path VARCHAR(250) NOT NULL,
    thumbnail BOOL NOT NULL,
    PRIMARY KEY(imageID),
    FOREIGN KEY(productID) REFERENCES Products(productID)
)