CREATE TABLE boxtype (
    id integer NOT NULL,
    name character varying(50) NOT NULL,
    CONSTRAINT pk_boxtype PRIMARY KEY (id)
);

INSERT INTO boxtype VALUES(1, 'Box (Man)');
INSERT INTO boxtype VALUES(2, 'Box (Woman)');
INSERT INTO boxtype VALUES(3, 'Box (Common)');
