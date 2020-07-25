CREATE TABLE box (
    idtype integer NOT NULL,
    "number" integer NOT NULL,
    card character varying(15),
    "time" timestamp without time zone NOT NULL DEFAULT now(),
    CONSTRAINT pk_uniqe_box_key UNIQUE (idtype, "number")
);

DO
$do$
    BEGIN
        FOR i IN 1..50 LOOP
                INSERT INTO box VALUES (1, i, '', now());
                INSERT INTO box VALUES (2, i, '', now());
                INSERT INTO box VALUES (3, i, '', now());
            END LOOP;
    END
$do$;
