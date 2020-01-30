CREATE TABLE userinfo (
  created timestamp without time zone NOT NULL DEFAULT now(),
  name character varying(25) NOT NULL,
  surname character varying(25) NOT NULL,
  phone character varying(20) NOT NULL,
  email character varying(50) NOT NULL,
  card character varying(15),
  balance integer NOT NULL DEFAULT 0,
  sex boolean,
  birthday timestamp without time zone,
  CONSTRAINT pk_user_info PRIMARY KEY (created),
  CONSTRAINT pk_uniqe_ui_key UNIQUE (email, phone),
  CONSTRAINT uk_userinfo_card UNIQUE (card)
);
