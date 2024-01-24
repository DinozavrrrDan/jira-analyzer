CREATE TABLE project (id serial PRIMARY KEY, title TEXT);

CREATE TABLE author (id serial PRIMARY KEY, name TEXT);

CREATE TABLE issues (
    id serial PRIMARY KEY,
    projectId INT NOT NULL,
    authorId INT NOT NULL,
    assigneeId INT NOT NULL,
    key TEXT,
    summary TEXT,
    description TEXT,
    type TEXT,
    priority TEXT,
    status TEXT,
    createdTime TIMESTAMP WITHOUT TIME ZONE,
    closedTime TIMESTAMP WITHOUT TIME ZONE,
    updatedTime TIMESTAMP WITHOUT TIME ZONE,
    timeSpent INT
);

GRANT ALL PRIVILEGES ON ALL TABLES IN SCHEMA public TO pguser;
GRANT USAGE, SELECT ON ALL SEQUENCES IN SCHEMA public TO pguser;
