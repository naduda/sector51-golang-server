CREATE TABLE usersecurity (
	password character varying(100) NOT NULL,
	roles character varying(100) NOT NULL DEFAULT 'USER',
	accountNonExpired boolean DEFAULT true,
	accountNonLocked boolean DEFAULT true,
	credentialsNonExpired boolean DEFAULT true,
	enabled boolean DEFAULT true,
	attempts integer NOT NULL DEFAULT 0,
	lastmodified timestamp without time zone NOT NULL DEFAULT now(),
	created timestamp without time zone NOT NULL DEFAULT now(),
	CONSTRAINT pk_user_security PRIMARY KEY (created)
);
