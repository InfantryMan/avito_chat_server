CREATE DATABASE db_chat;

CREATE USER chat_admin WITH password 'chat_password';

ALTER DATABASE db_chat OWNER TO chat_admin;

\connect db_chat

CREATE EXTENSION IF NOT EXISTS citext WITH SCHEMA public;

SET ROLE chat_admin;

CREATE TABLE "User"
(
    id         SERIAL PRIMARY KEY                          NOT NULL,
    username   CITEXT UNIQUE                               NOT NULL,
    created_at timestamptz DEFAULT transaction_timestamp() NOT NULL
);

CREATE TABLE "Chat"
(
    id         SERIAL PRIMARY KEY                          NOT NULL,
    name       CITEXT UNIQUE                               NOT NULL,
    created_at timestamptz DEFAULT transaction_timestamp() NOT NULL
);

CREATE TABLE "Chat_User"
(
    user_id INTEGER REFERENCES "User" (id) NOT NULL,
    chat_id INTEGER REFERENCES "Chat" (id) NOT NULL,
    UNIQUE (user_id, chat_id)
);

CREATE TABLE "Message"
(
    id         SERIAL PRIMARY KEY                          NOT NULL,
    text       text                                        NOT NULL,
    created_at timestamptz DEFAULT transaction_timestamp() NOT NULL,
    chat_id    INTEGER REFERENCES "Chat" (id)              NOT NULL,
    author_id  INTEGER REFERENCES "User" (id)              NOT NULL
);
