DROP TABLE if exists coupons;
DROP TABLE if exists brands;

CREATE TABLE brands
(
    ID serial NOT NULL PRIMARY KEY,
    name char(255) NOT NULL,
    CONSTRAINT "UNIQUE_brands_name" UNIQUE (name)

);

CREATE TABLE coupons
(
    id serial NOT NULL PRIMARY KEY,
    name char(255) NOT NULL,
    value money NOT NULL,
    created_at_utc timestamp NOT NULL default (timezone('utc', now())),
    expiry_utc timestamp NOT NULL,
    brand_id integer NOT NULL,
    CONSTRAINT "FK_brands" FOREIGN KEY (brand_id)
        REFERENCES brands (id)
        ON UPDATE CASCADE
        ON DELETE RESTRICT
);

CREATE INDEX "IDX_FK_brands" on coupons(brand_id);
CREATE INDEX "IDX_Expiry" on coupons(expiry_utc);
CREATE INDEX "IDX_value" on coupons(value);
CREATE INDEX "IDX_name" on brands(lower(name));
