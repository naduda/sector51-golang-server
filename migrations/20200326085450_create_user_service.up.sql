CREATE TABLE user_service (
      idservice integer NOT NULL,
      iduser timestamp without time zone NOT NULL,
      dtbeg timestamp without time zone,
      dtend timestamp without time zone,
      value character varying(50),
      CONSTRAINT pk_user_service PRIMARY KEY (idservice, iduser)
);
