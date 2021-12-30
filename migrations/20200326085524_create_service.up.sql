CREATE TABLE service (
    id integer NOT NULL,
    name character varying(25),
    "desc" character varying(50),
    price integer NOT NULL DEFAULT 0,
    CONSTRAINT pk_service PRIMARY KEY (id)
);
INSERT INTO service VALUES(0, 'ABONEMENT', '-', 700);
INSERT INTO service VALUES(1, 'TRAINER', '-', 10);
INSERT INTO service VALUES(2, 'BOX', '-', 50);
INSERT INTO service VALUES(3, 'ABONEMENT (Morning)', '-', 550);
INSERT INTO service VALUES(4, 'ABONEMENT (Evening)', '-', 600);
INSERT INTO service VALUES(5, 'ABONEMENT 3', '-', 1900);
INSERT INTO service VALUES(6, 'ABONEMENT (Morning) 3', '-', 1450);
INSERT INTO service VALUES(7, 'ABONEMENT (Evening) 3', '-', 1600);
INSERT INTO service VALUES(8, 'ABONEMENT 6', '-', 3150);
INSERT INTO service VALUES(9, 'ABONEMENT (Morning) 6', '-', 2400);
INSERT INTO service VALUES(10, 'ABONEMENT (Evening) 6', '-', 2800);
INSERT INTO service VALUES(11, 'ABONEMENT 12', '-', 4600);
INSERT INTO service VALUES(12, 'ABONEMENT (Morning) 12', '-', 4000);
INSERT INTO service VALUES(13, 'ABONEMENT (Evening) 12', '-', 4400);
INSERT INTO service VALUES(14, 'TWELVE', '-', 500);
