create table user_info
(
    user_no       serial not null
        constraint user_info_user_no_pk
            primary key,
    user_email    text   not null,
    user_name     text   not null,
    created_time  timestamp default now(),
    modified_time timestamp default now(),
    is_enabled    smallint  default 1,
    user_pw_hash  varchar(255)
);

alter table user_info
    owner to postgres;

create unique index user_info_email_uindex
    on user_info (user_email);

create unique index user_info_user_no_uindex
    on user_info (user_no);

create table smtp_info
(
    smtp_no       serial  not null
        constraint smtp_info_smtp_no_pk
            primary key,
    user_no       integer not null
        constraint smtp_info_user_info_user_no_fk
            references user_info
            on update cascade on delete cascade,
    smtp_host     text      default 'smtp.redteam.or.kr'::text,
    smtp_port     text      default '587'::text,
    protocol      text      default '1'::text,
    tls           text      default '1'::text,
    timeout       text      default '1000'::text,
    smtp_id       text,
    smtp_pw       text,
    created_time  timestamp default now(),
    modified_time timestamp default now()
);

alter table smtp_info
    owner to postgres;

create unique index smtp_info_smtp_no_uindex
    on smtp_info (smtp_no);

create table template_info
(
    tmp_no        serial not null
        constraint template_info_pk
            primary key,
    tmp_division  smallint,
    tmp_kind      smallint,
    file_info     smallint,
    tmp_name      text,
    mail_title    text,
    sender_name   text,
    download_type smallint,
    created_time  timestamp default now(),
    modified_time timestamp default now(),
    mail_content  text,
    user_no       integer
);

comment on column template_info.sender_name is 'email format';

alter table template_info
    owner to postgres;

create unique index template_info_tm_no_uindex
    on template_info (tmp_no);

create table tag_info
(
    tag_no        serial not null
        constraint tag_info_tag_no_pk
            primary key,
    tag_name      text,
    created_time  timestamp default now(),
    modified_time timestamp default now(),
    user_no       integer
        constraint tag_info_user_info_user_no_fk
            references user_info
            on update cascade on delete cascade
);

alter table tag_info
    owner to postgres;

create unique index tag_info_tag_no_uindex
    on tag_info (tag_no);

create table project_info
(
    p_no          serial not null
        constraint project_info_pk
            primary key,
    tml_no        integer
        constraint project_info_template_info_tmp_no_fk
            references template_info,
    tag_no        integer
        constraint project_info_tag_info_tag_no_fk
            references tag_info,
    p_name        text,
    p_description text,
    p_start_date  timestamp,
    p_end_date    timestamp,
    created_time  timestamp default now(),
    modified_time timestamp default now()
);

alter table project_info
    owner to postgres;

create unique index project_info_p_no_uindex
    on project_info (p_no);

create table tag_target_info
(
    target_no integer,
    tag_no    integer
        constraint tag_target_info_tag_info_tag_no_fk
            references tag_info,
    user_no   integer
        constraint tag_target_info_user_info_user_no_fk
            references user_info
            on update cascade on delete cascade
);

alter table tag_target_info
    owner to postgres;

create table target_info
(
    target_no       serial not null
        constraint target_info_target_no_pk
            primary key,
    target_name     text   not null,
    target_email    text   not null,
    target_phone    text,
    target_organize text,
    target_position text,
    created_time    timestamp default now(),
    modified_time   timestamp default now(),
    user_no         integer
        constraint target_info_user_info_user_no_fk
            references user_info
);

alter table target_info
    owner to postgres;


