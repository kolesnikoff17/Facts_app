BEGIN;
INSERT INTO facts(title, description) VALUES ('aboba', 'aboba_aboba');
INSERT INTO links(fact_id, link) VALUES (1, 'aboba');
INSERT INTO links(fact_id, link) VALUES (1, 'aboba_aboba');
END;

BEGIN;
INSERT INTO facts(title, description) VALUES ('aboba2', 'aboba_aboba2');
INSERT INTO links(fact_id, link) VALUES (2, 'aboba');
END;

INSERT INTO facts(title, description) VALUES ('aboba3', 'aboba_aboba3');

