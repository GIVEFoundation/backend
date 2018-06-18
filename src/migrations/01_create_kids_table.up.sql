create table kids
(
  id              character varying primary key not null,
  name            character varying not null,
  date_of_birth   date not null,
  parents_email   character varying[],
  students_photo  bytea,
  school_name     character varying,
  id_tag_name     character varying not null
);
